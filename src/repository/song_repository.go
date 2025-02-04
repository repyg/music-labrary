package repository

import (
	"log"
	"music-library/src/models"

	"gorm.io/gorm"
)

type SongRepository struct {
	db *gorm.DB
}

func NewSongRepository(db *gorm.DB) *SongRepository {
	return &SongRepository{db: db}
}

func (r *SongRepository) FindSongs(filter models.SongFilter) ([]models.Song, error) {
	var songs []models.Song
	query := r.db.Model(&models.Song{})

	if filter.Group != "" {
		query = query.Where("\"group\" ILIKE ?", "%"+filter.Group+"%")
	}
	if filter.Title != "" {
		query = query.Where("title ILIKE ?", "%"+filter.Title+"%")
	}
	if filter.ReleaseDate != "" {
		query = query.Where("release_date = ?", filter.ReleaseDate)
	}

	// Пагинация
	offset := (filter.Page - 1) * filter.PageSize
	query = query.Offset(offset).Limit(filter.PageSize)

	log.Printf("Executing query with filter: %+v", filter)
	if err := query.Find(&songs).Error; err != nil {
		return nil, err
	}

	return songs, nil
}

func (r *SongRepository) CreateSong(song models.Song) (*models.Song, error) {
	log.Printf("Creating song in database: %+v", song)
	if err := r.db.Create(&song).Error; err != nil {
		return nil, err
	}
	return &song, nil
}

func (r *SongRepository) DeleteSong(id uint) error {
	log.Printf("Deleting song with ID: %d", id)
	return r.db.Delete(&models.Song{}, id).Error
}

func (r *SongRepository) UpdateSong(song models.Song) (*models.Song, error) {
	log.Printf("Updating song with ID: %d", song.ID)
	if err := r.db.Save(&song).Error; err != nil {
		return nil, err
	}
	return &song, nil
}

func (r *SongRepository) FindById(id uint) (*models.Song, error) {
	var song models.Song
	if err := r.db.First(&song, id).Error; err != nil {
		return nil, err
	}
	return &song, nil
}
