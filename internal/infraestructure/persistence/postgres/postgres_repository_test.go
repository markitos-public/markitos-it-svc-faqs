package postgres_test

import (
	"testing"

	"markitos-it-svc-faqs/internal/domain/models"
	"markitos-it-svc-faqs/internal/domain/types"
	postgresrepo "markitos-it-svc-faqs/internal/infraestructure/persistence/postgres"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pgdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dsn = "host=localhost port=5431 user=postgres password=postgres dbname=faqs_test sslmode=disable"

func TestPostgresRepository_SaveAndGet(t *testing.T) {
	db, err := gorm.Open(pgdriver.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)

	repo := postgresrepo.NewPostgresFaqRepository(db)

	id, _ := types.NewID()
	title, _ := types.NewTitle("Test Title")
	content, _ := types.NewContent("This is a test content with more than ten characters.")
	tags, _ := types.NewTags([]string{"test", "go"})

	faq := &models.Faq{
		ID:      id,
		Title:   title,
		Content: content,
		Tags:    tags,
	}

	t.Cleanup(func() {
		err := db.Exec("DELETE FROM faqs WHERE id = ?", faq.ID.Value()).Error
		require.NoError(t, err)
		require.NoError(t, sqlDB.Close())
	})

	err = repo.Save(faq)
	require.NoError(t, err)

	retrievedFaq, err := repo.Get(id)
	require.NoError(t, err)
	assert.Equal(t, faq.ID.Value(), retrievedFaq.ID.Value())
	assert.Equal(t, faq.Title.Value(), retrievedFaq.Title.Value())
	assert.Equal(t, faq.Content.Value(), retrievedFaq.Content.Value())
	assert.Equal(t, faq.Tags.Value(), retrievedFaq.Tags.Value())
}

func TestPostgresRepository_Save_NilFaq(t *testing.T) {
	db, err := gorm.Open(pgdriver.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	repo := postgresrepo.NewPostgresFaqRepository(db)

	err = repo.Save(nil)
	assert.Error(t, err)
	assert.Equal(t, "faq is nil", err.Error())
}

func TestPostgresRepository_Save_IncompleteFaq(t *testing.T) {
	db, err := gorm.Open(pgdriver.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	repo := postgresrepo.NewPostgresFaqRepository(db)

	id, _ := types.NewID()
	title, _ := types.NewTitle("Test Title")
	content, _ := types.NewContent("This is a test content with more than ten characters.")

	faq := &models.Faq{
		ID:      id,
		Title:   title,
		Content: content,
		Tags:    nil, // Tags is nil
	}

	err = repo.Save(faq)
	assert.Error(t, err)
	assert.Equal(t, "faq is incomplete", err.Error())
}

func TestPostgresRepository_Get_NonExistentFaq(t *testing.T) {
	db, err := gorm.Open(pgdriver.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	repo := postgresrepo.NewPostgresFaqRepository(db)

	nonExistentID, _ := types.NewID()

	faq, err := repo.Get(nonExistentID)
	assert.Error(t, err)
	assert.Nil(t, faq)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestPostgresRepository_DeleteExistingFaq(t *testing.T) {
	db, err := gorm.Open(pgdriver.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	repo := postgresrepo.NewPostgresFaqRepository(db)

	id, _ := types.NewID()
	title, _ := types.NewTitle("Test Title")
	content, _ := types.NewContent("This is a test content with more than ten characters.")
	tags, _ := types.NewTags([]string{"test", "go"})

	faq := &models.Faq{
		ID:      id,
		Title:   title,
		Content: content,
		Tags:    tags,
	}

	err = repo.Save(faq)
	require.NoError(t, err)

	err = repo.Delete(id)
	require.NoError(t, err)

	deletedFaq, err := repo.Get(id)
	assert.Error(t, err)
	assert.Nil(t, deletedFaq)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestPostgresRepository_DeleteNonExistentFaq(t *testing.T) {
	db, err := gorm.Open(pgdriver.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	repo := postgresrepo.NewPostgresFaqRepository(db)

	nonExistentID, _ := types.NewID()

	err = repo.Delete(nonExistentID)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestPostgresRepository_Delete_NilID(t *testing.T) {
	db, err := gorm.Open(pgdriver.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	repo := postgresrepo.NewPostgresFaqRepository(db)

	err = repo.Delete(nil)
	assert.Error(t, err)
	assert.Equal(t, "faq id is nil", err.Error())
}

func TestPostgresRepository_UpdateExistingFaq(t *testing.T) {
	db, err := gorm.Open(pgdriver.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	repo := postgresrepo.NewPostgresFaqRepository(db)

	id, _ := types.NewID()
	title, _ := types.NewTitle("Test Title")
	content, _ := types.NewContent("This is a test content with more than ten characters.")
	tags, _ := types.NewTags([]string{"test", "go"})

	faq := &models.Faq{
		ID:      id,
		Title:   title,
		Content: content,
		Tags:    tags,
	}

	err = repo.Save(faq)
	require.NoError(t, err)

	newTitle, _ := types.NewTitle("Updated Title")
	newContent, _ := types.NewContent("This is updated content with more than ten characters.")
	newTags, _ := types.NewTags([]string{"updated", "go"})

	faq.Title = newTitle
	faq.Content = newContent
	faq.Tags = newTags

	err = repo.Update(faq)
	require.NoError(t, err)

	updatedFaq, err := repo.Get(id)
	require.NoError(t, err)
	assert.Equal(t, faq.ID.Value(), updatedFaq.ID.Value())
	assert.Equal(t, faq.Title.Value(), updatedFaq.Title.Value())
	assert.Equal(t, faq.Content.Value(), updatedFaq.Content.Value())
	assert.Equal(t, faq.Tags.Value(), updatedFaq.Tags.Value())
}

func TestPostgresRepository_List(t *testing.T) {
	db, err := gorm.Open(pgdriver.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)

	repo := postgresrepo.NewPostgresFaqRepository(db)

	id1, _ := types.NewID()
	title1, _ := types.NewTitle("First Test Title")
	content1, _ := types.NewContent("This is the first test content with enough characters.")
	tags1, _ := types.NewTags([]string{"first", "test"})

	faq1 := &models.Faq{
		ID:      id1,
		Title:   title1,
		Content: content1,
		Tags:    tags1,
	}

	id2, _ := types.NewID()
	title2, _ := types.NewTitle("Second Test Title")
	content2, _ := types.NewContent("This is the second test content with enough characters.")
	tags2, _ := types.NewTags([]string{"second", "test"})

	faq2 := &models.Faq{
		ID:      id2,
		Title:   title2,
		Content: content2,
		Tags:    tags2,
	}

	t.Cleanup(func() {
		err := db.Exec("DELETE FROM faqs WHERE id IN (?, ?)", faq1.ID.Value(), faq2.ID.Value()).Error
		require.NoError(t, err)
		require.NoError(t, sqlDB.Close())
	})

	err = repo.Save(faq1)
	require.NoError(t, err)

	err = repo.Save(faq2)
	require.NoError(t, err)

	faqs, err := repo.List()
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(faqs), 2)

	found1 := false
	found2 := false
	for _, a := range faqs {
		if a.ID.Value() == faq1.ID.Value() {
			found1 = true
			assert.Equal(t, faq1.Title.Value(), a.Title.Value())
			assert.Equal(t, faq1.Content.Value(), a.Content.Value())
			assert.Equal(t, faq1.Tags.Value(), a.Tags.Value())
		}
		if a.ID.Value() == faq2.ID.Value() {
			found2 = true
			assert.Equal(t, faq2.Title.Value(), a.Title.Value())
			assert.Equal(t, faq2.Content.Value(), a.Content.Value())
			assert.Equal(t, faq2.Tags.Value(), a.Tags.Value())
		}
	}

	assert.True(t, found1, "faq1 should be present in the list")
	assert.True(t, found2, "faq2 should be present in the list")
}
