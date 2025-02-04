package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"music-library/src/config"
	"music-library/src/models"
	"music-library/src/repository"
)

type SongService struct {
	repo        *repository.SongRepository
	externalURL string
	config      *config.Config
}

func NewSongService(repo *repository.SongRepository, externalURL string, config *config.Config) *SongService {
	return &SongService{
		repo:        repo,
		externalURL: externalURL,
		config:      config,
	}
}

func (s *SongService) GetSongs(filter models.SongFilter) ([]models.Song, error) {
	log.Printf("Getting songs with filter: %+v", filter)
	return s.repo.FindSongs(filter)
}

func (s *SongService) AddSong(song models.Song) (*models.Song, error) {
	log.Printf("Debug: Requesting external API for song details: %s - %s", song.Group, song.Title)

	url := fmt.Sprintf("%s/info?group=%s&song=%s", s.config.ExternalAPIURL, url.QueryEscape(song.Group), url.QueryEscape(song.Title))
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error: Failed to request external API: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error: External API returned status %d", resp.StatusCode)
		return nil, fmt.Errorf("external API returned status %d", resp.StatusCode)
	}

	var songDetail models.SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
		log.Printf("Error: Failed to decode API response: %v", err)
		return nil, err
	}

	log.Printf("Info: Successfully retrieved song details from external API")

	song.ReleaseDate = songDetail.ReleaseDate
	song.Text = songDetail.Text
	song.Link = songDetail.Link

	return s.repo.CreateSong(song)
}

func (s *SongService) DeleteSong(id uint) error {
	return s.repo.DeleteSong(id)
}

func (s *SongService) UpdateSong(song models.Song) (*models.Song, error) {
	return s.repo.UpdateSong(song)
}

func (s *SongService) GetSongById(id uint) (*models.Song, error) {
	return s.repo.FindById(id)
}

func (s *SongService) GetSongLyrics(id uint, page, versesPerPage int) (*models.Lyrics, error) {
	song, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	verses := strings.Split(song.Text, "\n\n")
	totalPages := (len(verses) + versesPerPage - 1) / versesPerPage

	start := (page - 1) * versesPerPage
	end := start + versesPerPage
	if end > len(verses) {
		end = len(verses)
	}

	return &models.Lyrics{
		SongID:        song.ID,
		Title:         song.Title,
		Group:         song.Group,
		Verses:        verses[start:end],
		CurrentPage:   page,
		TotalPages:    totalPages,
		VersesPerPage: versesPerPage,
	}, nil
}

func (s *SongService) fetchSongDetails(group, title string) (*models.SongDetail, error) {
	url := fmt.Sprintf("%s/info?group=%s&song=%s", s.externalURL, group, title)

	log.Printf("Fetching song details from: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external API returned status: %d", resp.StatusCode)
	}

	var details models.SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return nil, err
	}

	return &details, nil
}

type AddSongRequest struct {
	Group string `json:"group" binding:"required"`
	Song  string `json:"song" binding:"required"`
}
