package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"
	"parcel"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
	"fmt"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange = rand.New(randSource)
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

	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора *оставляю комменты для себя, чтобы не запутаться, т.к. проект еще буду править*
	testParcel := getTestParcel()
	id, err := Add(testParcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	// получите только что добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что значения всех полей в полученном объекте совпадают со значениями полей в переменной parcel
	tmpParcel, er := Get(id)
	require.NoError(t, er)

	assert.Equal(t, testParcel.Client, tmpParcel.Client)
	assert.Equal(t, testParcel.Status, tmpParcel.Status)
	assert.Equal(t, testParcel.Address, tmpParcel.Address)
	assert.Equal(t, testParcel.CreatedAt, tmpParcel.CreatedAt)

	// удалите добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что посылку больше нельзя получить из БД
	errr := delete(id)
	require.NoError(t, errr)

	_, errrr := Get(id)
	require.Equal(t, errrr, sql.ErrNoRows)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	testParcel := getTestParcel()
	id, err := Add(testParcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	// обновите адрес, убедитесь в отсутствии ошибки
	newAddress := "new test address"
	er := SetAddress(id, newAddress)
	require.NoError(t, er)

	// получите добавленную посылку и убедитесь, что адрес обновился
	tmpParcel, errr := Get(id)
	require.NoError(t, errr)
	require.Equal(t, tmpParcel.Address, newAddress)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	testParcel := getTestParcel()
	id, err := Add(testParcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	// обновите статус, убедитесь в отсутствии ошибки
	er := SetStatus(id, ParcelStatusSent)
	require.NoError(t, er)

	// получите добавленную посылку и убедитесь, что статус обновился
	tmpParcel, errr := Get(id)
	require.NoError(t, errr)
	require.Equal(t, tmpParcel.Status, ParcelStatusSent)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	for i := 0; i < len(parcels); i++ {
		id, err := Add(parcels[i]) // добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
		require.NoError(t, err)
		require.NotEmpty(t, id)

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	storedParcels, err := GetByClient(client)	// получите список посылок по идентификатору клиента, сохранённого в переменной client
	require.NoError(t, err)	// убедитесь в отсутствии ошибки
	require.Equal(t, len(storedParcels), len(parcelMap)) 	// убедитесь, что количество полученных посылок совпадает с количеством добавленных


	fmt.Println(storedParcels)
	fmt.Println(parcelMap)
	for _, parcel := range storedParcels {
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		// убедитесь, что все посылки из storedParcels есть в parcelMap
		// убедитесь, что значения полей полученных посылок заполнены верно
		require.Equal(t.Parcel.Number, parcelMap[t.Parcel.Number].Number)
		require.Equal(t.Parcel.Client, parcelMap[t.Parcel.Number].Client)
		require.Equal(t.Parcel.Status, parcelMap[t.Parcel.Number].Status)
		require.Equal(t.Parcel.Address, parcelMap[t.Parcel.Number].Address)
		require.Equal(t.Parcel.CreatedAt, parcelMap[t.Parcel.Number].CreatedAt)
	}
}
