package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)

	testParcel := getTestParcel()
	id, err := store.Add(testParcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	tmpParcel, err := store.Get(id)
	require.NoError(t, err)

	assert.Equal(t, testParcel.Client, tmpParcel.Client)
	assert.Equal(t, testParcel.Status, tmpParcel.Status)
	assert.Equal(t, testParcel.Address, tmpParcel.Address)
	assert.Equal(t, testParcel.CreatedAt, tmpParcel.CreatedAt)

	err = store.Delete(id)
	require.NoError(t, err)

	_, err = store.Get(id)
	require.Equal(t, err, sql.ErrNoRows)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)

	testParcel := getTestParcel()
	id, err := store.Add(testParcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	tmpParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, tmpParcel.Address, newAddress)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)

	testParcel := getTestParcel()
	id, err := store.Add(testParcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	err = store.SetStatus(id, ParcelStatusSent)
	require.NoError(t, err)

	tmpParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, tmpParcel.Status, ParcelStatusSent)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.NotEmpty(t, id)

		parcels[i].Number = id

		parcelMap[id] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	require.Equal(t, len(storedParcels), len(parcelMap))

	for _, parcel := range storedParcels {
		require.Equal(t, parcel.Number, parcelMap[parcel.Number].Number)
		require.Equal(t, parcel.Client, parcelMap[parcel.Number].Client)
		require.Equal(t, parcel.Status, parcelMap[parcel.Number].Status)
		require.Equal(t, parcel.Address, parcelMap[parcel.Number].Address)
		require.Equal(t, parcel.CreatedAt, parcelMap[parcel.Number].CreatedAt)
	}
}
