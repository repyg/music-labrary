package api

import (
	"log"
	"strconv"

	"music-library/src/models"
	"music-library/src/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.SongService
}

func NewHandler(service *service.SongService) *Handler {
	return &Handler{service: service}
}

// GetSongs godoc
// @Summary Получение списка песен с фильтрацией и пагинацией
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string false "Группа"
// @Param title query string false "Название песни"
// @Param page query int false "Номер страницы"
// @Param size query int false "Размер страницы"
// @Success 200 {array} models.Song
// @Router /songs [get]
func (h *Handler) GetSongs(c *gin.Context) {
	group := c.Query("group")
	title := c.Query("title")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	filter := models.SongFilter{
		Group:    group,
		Title:    title,
		Page:     page,
		PageSize: pageSize,
	}

	songs, err := h.service.GetSongs(filter)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get songs"})
		return
	}

	c.JSON(200, songs)
}

// AddSongRequest представляет формат входящего JSON
type AddSongRequest struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

// AddSong godoc
// @Summary Добавление новой песни
// @Tags songs
// @Accept json
// @Produce json
// @Param song body AddSongRequest true "Данные песни"
// @Success 201 {object} models.Song
// @Router /songs [post]
func (h *Handler) AddSong(c *gin.Context) {
	var req AddSongRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Debug: Invalid request body: %v", err)
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	song := models.Song{
		Group: req.Group,
		Title: req.Song,
	}

	newSong, err := h.service.AddSong(song)
	if err != nil {
		log.Printf("Error: Failed to add song: %v", err)
		c.JSON(500, gin.H{"error": "Failed to add song"})
		return
	}

	log.Printf("Info: Successfully added song: %s by %s", req.Song, req.Group)
	c.JSON(201, newSong)
}

// DeleteSong godoc
// @Summary Удаление песни
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Success 204
// @Router /songs/{id} [delete]
func (h *Handler) DeleteSong(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.service.DeleteSong(uint(id)); err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete song"})
		return
	}

	c.Status(204)
}

// UpdateSong godoc
// @Summary Изменение данных песни
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param song body models.Song true "Данные песни"
// @Success 200 {object} models.Song
// @Router /songs/{id} [put]
func (h *Handler) UpdateSong(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var song models.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	song.ID = uint(id)

	updatedSong, err := h.service.UpdateSong(song)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update song"})
		return
	}

	c.JSON(200, updatedSong)
}

// GetSongLyrics godoc
// @Summary Получение текста песни с пагинацией по куплетам
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param verse query int false "Номер куплета"
// @Success 200 {string} string
// @Router /songs/{id}/lyrics [get]
func (h *Handler) GetSongLyrics(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	verse, _ := strconv.Atoi(c.DefaultQuery("verse", "1"))
	versesPerPage, _ := strconv.Atoi(c.DefaultQuery("versesPerPage", "5"))

	lyrics, err := h.service.GetSongLyrics(uint(id), verse, versesPerPage)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get song lyrics"})
		return
	}

	c.JSON(200, lyrics)
}

// GetSongById godoc
// @Summary Получение песни по ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Success 200 {object} models.Song
// @Router /songs/{id} [get]
func (h *Handler) GetSongById(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	song, err := h.service.GetSongById(uint(id))
	if err != nil {
		log.Printf("Error getting song by ID: %v", err)
		c.JSON(404, gin.H{"error": "Song not found"})
		return
	}

	c.JSON(200, song)
}
