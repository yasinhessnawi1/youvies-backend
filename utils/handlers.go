package utils

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"youvies-backend/models"
)

// ExtractQuality extracts the quality of the torrent from its name.
func ExtractQuality(name string) string {
	name = strings.ToLower(name)
	switch {
	case strings.Contains(name, "1080p"):
		return "1080p"
	case strings.Contains(name, "720p") || strings.Contains(name, "dvd") || strings.Contains(name, "dvdrip") ||
		strings.Contains(name, "web-dl") || strings.Contains(name, "webdl") || strings.Contains(name, "webrip"):
		return "720p"
	case strings.Contains(name, "480p"):
		return "480p"
	case strings.Contains(name, "4k"), strings.Contains(name, "2160p"):
		return "4k"
	default:
		return "unknown"
	}
}

// SortTorrentsBySeeders sorts torrents by the number of seeders.
func SortTorrentsBySeeders(torrents []models.Torrent) {
	sort.SliceStable(torrents, func(i, j int) bool {
		return torrents[i].Seeders > torrents[j].Seeders
	})
}

// isUnrelatedContent filters out unrelated torrents based on their name.
func isUnrelatedContent(name string) bool {
	unrelatedKeywords := []string{
		"book", "guide", "soundtrack", "companion", "album", "cookbook", "unofficial",
	}
	for _, keyword := range unrelatedKeywords {
		if strings.Contains(strings.ToLower(name), keyword) {
			return true
		}
	}
	return false
}

// GetEpisodeNumberFromTorrentName extracts the episode number from the torrent name.
func GetEpisodeNumberFromTorrentName(name string) int {
	episodePatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)\bS(\d{1,2})E(\d{1,3})\b`),                   // S01E01, S1E01
		regexp.MustCompile(`(?i)\b(\d{1,2})x(\d{1,3})\b`),                    // 1x01, 01x01
		regexp.MustCompile(`(?i)season\s?(\d{1,2})\s?episode\s?(\d{1,3})\b`), // season 1 episode 1
		regexp.MustCompile(`(?i)episode\s?(\d{1,3})\b`),                      // episode 1
		regexp.MustCompile(`(?i)[Ee]p?(\d{1,3})\b`),                          // E01, Ep01, ep1
		regexp.MustCompile(`(?i) - (\d{1,3})\b`),                             // - 01
		regexp.MustCompile(`(?i)\[(\d{1,3})\]\b`),                            // [01]
		regexp.MustCompile(`(?i)\s(\d{1,3})\b`),                              // 01 (space followed by number)
		regexp.MustCompile(`(?i)第(\d{1,3})話`),                                // 第01話 (Japanese)
		regexp.MustCompile(`(?i)\b(\d{1,3})話\b`),                             // 01話 (Japanese)
		regexp.MustCompile(`(?i)\b(\d{1,3})화\b`),                             // 01화 (Korean)
		regexp.MustCompile(`(?i) - (\d{1,3})화\b`),                            // - 01화 (Korean)
		regexp.MustCompile(`(?i)ep(\d{1,3})`),                                // ep01
		regexp.MustCompile(`(?i)ep\s(\d{1,3})`),                              // ep 01
		regexp.MustCompile(`(?i)ep\.(\d{1,3})`),                              // ep.01
		regexp.MustCompile(`(?i)episode\.(\d{1,3})`),                         // episode.01
		regexp.MustCompile(`(?i)e(\d{1,3})`),                                 // e01
		regexp.MustCompile(`(?i)\b(\d{1,3})\b`),                              // 01
		regexp.MustCompile(`(?i)[._](\d{1,3})[._]`),                          // .01., _01_
	}

	for _, pattern := range episodePatterns {
		matches := pattern.FindStringSubmatch(name)
		if len(matches) >= 2 {
			episodeNum, err := strconv.Atoi(matches[1])
			if err == nil {
				return episodeNum
			}
		}
	}
	return 0
}

// ExtractSeasonAndEpisode extracts the season and episode numbers from the torrent name.
func ExtractSeasonAndEpisode(name string) (int, int, error) {
	// Regex patterns to capture season and episode numbers
	patterns := []string{
		`S(\d{1,2})E(\d{1,2})`,                        // Standard format S01E01
		`(\d{1,2})x(\d{1,2})`,                         // Alternate format 1x01
		`Season[ _](\d{1,2})[ _]Episode[ _](\d{1,2})`, // Format Season 1 Episode 1
		`Episode[ _](\d{1,2})`,                        // Format Episode 1
		`[Ee]p?(\d{1,2})`,                             // Format E01, Ep01, ep1
		` - (\d{1,2})`,                                // Format - 01
		`\[(\d{1,2})\]`,                               // Format [01]
		`\s(\d{1,2})`,                                 // Format 01 (space followed by number)
		`第(\d{1,2})話`,                                 // Format 第01話 (Japanese)
		`(\d{1,2})話`,                                  // Format 01話 (Japanese)
		`(\d{1,2})화`,                                  // Format 01화 (Korean)
		` - (\d{1,2})화`,                               // Format - 01화 (Korean)
		`ep(\d{1,2})`,                                 // Format ep01
		`ep\s(\d{1,2})`,                               // Format ep 01
		`episode\.(\d{1,2})`,                          // Format episode.01
		`ep\.(\d{1,2})`,                               // Format ep.01
		`e(\d{1,2})`,                                  // Format e01
		`(\d{1,2})`,                                   // Format 01
		`[._](\d{1,2})[._]`,                           // Format .01., _01_
		`(\d{1,2})\D+(\d{1,2})`,                       // Loose format 01-01 or 01x01
		`(\d{1,2})\D*(\d{1,2})`,                       // Another loose format
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(name)
		if len(matches) >= 3 {
			season, err := strconv.Atoi(matches[1])
			if err != nil {
				continue
			}
			episode, err := strconv.Atoi(matches[2])
			if err != nil {
				continue
			}
			return season, episode, nil
		}
	}

	return 0, 0, fmt.Errorf("no season or episode number found in: %s", name)
}

// CategorizeAnimeTorrentsBySeasonsAndEpisodes categorizes anime torrents by seasons based on year and episodes.
func CategorizeAnimeTorrentsBySeasonsAndEpisodes(torrents []models.Torrent, episodeCount int) (map[int]models.Season, map[string][]models.Torrent) {
	seasons := make(map[int]models.Season)
	fullContent := make(map[string][]models.Torrent)
	seasonNum := 1

	for _, torrent := range torrents {
		if isUnrelatedContent(torrent.Name) {
			continue
		}

		// Check for full season or complete series torrents
		if isFullSeasonTorrent(torrent.Name) {
			fullContent[torrent.Name] = append(fullContent[torrent.Name], torrent)
			continue
		}

		season, exists := seasons[seasonNum]
		if !exists {
			season = models.Season{Episodes: make(map[int]models.Episode)}
		}

		episodeNum := GetEpisodeNumberFromTorrentName(torrent.Name)
		if episodeNum == 0 {
			log.Printf("error parsing episode number: %s", torrent.Name)
			continue
		}

		episode, exists := season.Episodes[episodeNum]
		if !exists {
			episode = models.Episode{Torrents: make(map[string][]models.Torrent)}
		}

		quality := ExtractQuality(torrent.Name)
		episode.Torrents[quality] = append(episode.Torrents[quality], torrent)

		season.Episodes[episodeNum] = episode
		seasons[seasonNum] = season
	}

	// Sort torrents by seeders within each episode's quality category
	for _, season := range seasons {
		for _, episode := range season.Episodes {
			for quality := range episode.Torrents {
				SortTorrentsBySeeders(episode.Torrents[quality])
			}
		}
	}

	return seasons, fullContent
}

// isFullSeasonTorrent checks if a torrent is for a full season or complete series.
func isFullSeasonTorrent(name string) bool {
	fullSeasonPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)\b(complete|full)\b`),                                    // contains "complete" or "full"
		regexp.MustCompile(`(?i)\b(season\s?\d{1,2}|s\d{1,2})\b`),                        // Season 1, S01
		regexp.MustCompile(`(?i)\b(season\s?\d{1,2}|s\d{1,2})\s?(complete|full)`),        // Season 1 complete, S01 full
		regexp.MustCompile(`(?i)\b(complete|full)\s?(season\s?\d{1,2}|s\d{1,2})`),        // complete season 1, full S01
		regexp.MustCompile(`(?i)\b(season\s?\d{1,2}|s\d{1,2})\s?(series|complete|full)`), // Season 1 series, S01 complete, S01 full
		regexp.MustCompile(`(?i)\b(series|complete|full)\s?(season\s?\d{1,2}|s\d{1,2})`), // series Season 1, complete S01, full S01
		regexp.MustCompile(`(?i)\bseries\s?(complete|full)\b`),                           // series complete, series full
		regexp.MustCompile(`(?i)\b(complete|full)\s?series\b`),                           // complete series, full series

	}

	for _, pattern := range fullSeasonPatterns {
		if pattern.MatchString(name) {
			return true
		}
	}

	return false
}

// FetchMissingTorrents fetches missing torrents for a show based on the season and episode numbers.
func FetchMissingTorrents(title string, existingTorrents []models.Torrent, seasonsInfo []models.SeasonInfo) ([]models.Torrent, error) {
	var missingTorrents []models.Torrent
	// Determine which episodes are missing
	missingEpisodes := make(map[int][]int)
	for _, seasonInfo := range seasonsInfo {
		if seasonInfo.SeasonNumber == 0 {
			continue // Skip season 0
		}
		for episodeNum := 1; episodeNum <= seasonInfo.EpisodeCount; episodeNum++ {
			missingEpisodes[seasonInfo.SeasonNumber] = append(missingEpisodes[seasonInfo.SeasonNumber], episodeNum)
		}
	}

	// Remove found episodes
	for _, torrent := range existingTorrents {
		seasonNum, episodeNum, err := ExtractSeasonAndEpisode(torrent.Name)
		if err != nil {
			log.Printf("error parsing season and episode number: %v", err)
			continue
		}
		episodes := missingEpisodes[seasonNum]
		for i, ep := range episodes {
			if ep == episodeNum {
				missingEpisodes[seasonNum] = append(episodes[:i], episodes[i+1:]...)
				break
			}
		}
	}

	// Fetch missing torrents
	for seasonNum, episodes := range missingEpisodes {
		for _, episodeNum := range episodes {
			query := fmt.Sprintf("%s S%02dE%02d", title, seasonNum, episodeNum)
			torrents, err := FetchTorrents(query)
			if err != nil {
				log.Printf("Failed to fetch torrents for %s: %v", query, err)
				continue
			}
			for _, torrent := range torrents {
				err := SaveMetadata(torrent.Magnet, torrent.Name)
				if err != nil {
					log.Printf("Failed to save torrent metadata for %s: %v", torrent.Name, err)
					continue
				}
			}
			missingTorrents = append(missingTorrents, torrents...)
		}
	}
	if len(missingTorrents) < len(missingEpisodes)-5 {
		return nil, fmt.Errorf("Error fetching all messing episodes: actual missing number %d, found %d ", len(missingEpisodes), len(missingTorrents))
	}

	return missingTorrents, nil
}

// CategorizeTorrentsBySeasonsAndEpisodes categorizes torrents by their respective seasons and episodes.
func CategorizeTorrentsBySeasonsAndEpisodes(torrents []models.Torrent) (map[int]models.Season, []models.Torrent) {
	seasons := make(map[int]models.Season)
	extra := make([]models.Torrent, 0)

	for _, torrent := range torrents {
		seasonNum, episodeNum, err := ExtractSeasonAndEpisode(torrent.Name)
		if err != nil {
			extra = append(extra, torrent)
			continue
		}
		if seasonNum == 0 {
			continue // Skip season 0
		}
		if _, exists := seasons[seasonNum]; !exists {
			seasons[seasonNum] = models.Season{
				Episodes: make(map[int]models.Episode),
			}
		}
		season := seasons[seasonNum]
		if _, exists := season.Episodes[episodeNum]; !exists {
			season.Episodes[episodeNum] = models.Episode{
				Torrents: make(map[string][]models.Torrent),
			}
		}
		episode := season.Episodes[episodeNum]
		quality := ExtractQuality(torrent.Name)
		episode.Torrents[quality] = append(episode.Torrents[quality], torrent)
		season.Episodes[episodeNum] = episode
		seasons[seasonNum] = season
	}
	// Sort torrents by seeders within each episode's quality category
	for _, season := range seasons {
		for _, episode := range season.Episodes {
			for quality := range episode.Torrents {
				SortTorrentsBySeeders(episode.Torrents[quality])
			}
		}
	}

	return seasons, extra
}

// CategorizeTorrentsByQuality categorizes torrents by their quality.
func CategorizeTorrentsByQuality(torrents []models.Torrent) map[string][]models.Torrent {
	categorized := make(map[string][]models.Torrent)
	for _, torrent := range torrents {
		quality := ExtractQuality(torrent.Name)
		categorized[quality] = append(categorized[quality], torrent)
	}
	for quality := range categorized {
		SortTorrentsBySeeders(categorized[quality])
	}
	return categorized
}
