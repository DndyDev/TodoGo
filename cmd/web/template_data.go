// templates
package main

import (
	"dandydev.com/todogo/pkg/models"
)

type templateData struct {
	Note  *models.Note
	Notes []*models.Note
}
