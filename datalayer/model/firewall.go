package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Firewall struct {
	ID           bson.ObjectID   `bson:"_id,omitempty"           json:"id"`
	Name         string          `bson:"name"                    json:"name"`
	Enabled      bool            `bson:"enabled"                 json:"enabled"`
	TrustHeaders []string        `bson:"trust_headers,omitempty" json:"trust_headers,omitzero"`
	TrustProxies []string        `bson:"trust_proxies,omitempty" json:"trust_proxies,omitzero"`
	Blacklist    []string        `bson:"blacklist,omitempty"     json:"blacklist,omitzero"` // IP 白名单（与地区模式二选一）
	Regions      FirewallRegions `bson:"regions,omitempty"       json:"regions,omitzero"`   // 地区白名单（与 IP 白名单模式二选一）
	CreatedAt    time.Time       `bson:"created_at,omitempty"    json:"created_at"`
	UpdatedAt    time.Time       `bson:"updated_at,omitempty"    json:"updated_at"`
}

type FirewallRegion struct {
	Country string   `bson:"country" json:"country" validate:"required,lte=100"`                      // 国家/地区
	Cities  []string `bson:"cities"  json:"cities"  validate:"lte=1000,unique,dive,required,lte=500"` // 城市
}

type FirewallRegions []*FirewallRegion

func (frs FirewallRegions) Table() *FirewallRegionTable {
	tables := make(map[string]map[string]struct{}, len(frs))
	for _, fr := range frs {
		country := fr.Country
		cities := tables[country]
		if cities == nil {
			cities = make(map[string]struct{}, len(fr.Cities))
		}
		for _, city := range fr.Cities {
			cities[city] = struct{}{}
		}
		tables[country] = cities
	}

	return &FirewallRegionTable{tables: tables}
}

// Merge 将重复的国家城市合并。
func (frs FirewallRegions) Merge() FirewallRegions {
	size := len(frs)
	countries := make(map[string]*FirewallRegion, size)
	uniq := make(map[string]map[string]struct{}, len(frs))
	results := make(FirewallRegions, 0, size)
	for _, fr := range frs {
		country := fr.Country
		region := countries[country]
		if region == nil {
			region = &FirewallRegion{Country: country}
			results = append(results, region)
			countries[country] = region
			uniq[country] = make(map[string]struct{}, len(fr.Cities))
		}
		cities := uniq[country]
		for _, city := range fr.Cities {
			if _, exists := cities[city]; exists {
				continue
			}
			cities[city] = struct{}{}
			region.Cities = append(region.Cities, city)
		}
	}

	return results
}

type FirewallRegionTable struct {
	tables map[string]map[string]struct{}
}

// Contains 判断一个国家+城市是否在内。
func (frt FirewallRegionTable) Contains(country, city string) bool {
	cities, exists := frt.tables[country]
	if !exists {
		return false
	} else if len(cities) == 0 {
		return true
	}
	_, exists = cities[city]

	return exists
}
