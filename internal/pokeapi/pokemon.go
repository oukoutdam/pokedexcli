package pokeapi

import (
	"net/url"
)

type Pokemon struct {
	ID                     int           `json:"id"`
	IsDefault              bool          `json:"is_default"`
	LocationAreaEncounters string        `json:"location_area_encounters"`
	Height                 int           `json:"height"`
	Weight                 int           `json:"weight"`
	BaseExperience         int           `json:"base_experience"`
	Name                   string        `json:"name"`
	Order                  int           `json:"order"`
	Species                NamedResource `json:"species"`
	Abilities              []struct {
		Ability  NamedResource `json:"ability"`
		IsHidden bool          `json:"is_hidden"`
		Slot     int           `json:"slot"`
	} `json:"abilities"`
	Cries struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Forms       []NamedResource `json:"forms"`
	GameIndices []struct {
		GameIndex int           `json:"game_index"`
		Version   NamedResource `json:"version"`
	} `json:"game_indices"`
	HeldItems []struct {
		Item           NamedResource `json:"item"`
		VersionDetails []struct {
			Rarity  int           `json:"rarity"`
			Version NamedResource `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	Moves []struct {
		Move                NamedResource `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int           `json:"level_learned_at"`
			MoveLearnMethod NamedResource `json:"move_learn_method"`
			Order           any           `json:"order"`
			VersionGroup    NamedResource `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Stats []struct {
		BaseStat int           `json:"base_stat"`
		Effort   int           `json:"effort"`
		Stat     NamedResource `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int           `json:"slot"`
		Type NamedResource `json:"type"`
	} `json:"types"`
}

func (c *Client) GetPokemon(pokemonName string) (Pokemon, error) {
	return getResource[Pokemon](c, baseURL+"/pokemon/"+url.PathEscape(pokemonName))
}
