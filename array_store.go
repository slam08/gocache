package cache

import (
	"errors"
	"strconv"
)

// ArrayStore is the representation of an array caching store
type ArrayStore struct {
	Client map[string]interface{}
	Prefix string
}

// Get a value from the store
func (as *ArrayStore) Get(key string) (interface{}, error) {
	value := as.Client[as.GetPrefix()+key]

	if value == nil {
		return "", nil
	}

	if IsStringNumeric(value.(string)) {
		floatValue, err := StringToFloat64(value.(string))

		if err != nil {
			return floatValue, err
		}

		if IsFloat(floatValue) {
			return floatValue, err
		}

		return int64(floatValue), err
	}

	return SimpleDecode(value.(string))
}

// Get a float value from the store
func (as *ArrayStore) GetFloat(key string) (float64, error) {
	value := as.Client[as.GetPrefix()+key]

	if value == nil || !IsStringNumeric(value.(string)) {
		return 0, errors.New("Invalid numeric value")
	}

	return StringToFloat64(value.(string))
}

// Get an int value from the store
func (as *ArrayStore) GetInt(key string) (int64, error) {
	value := as.Client[as.GetPrefix()+key]

	if value == nil || !IsStringNumeric(value.(string)) {
		return 0, errors.New("Invalid numeric value")
	}

	val, err := StringToFloat64(value.(string))

	return int64(val), err
}

// Increment an integer counter by a given value
func (as *ArrayStore) Increment(key string, value int64) (int64, error) {
	val := as.Client[as.GetPrefix()+key]

	if val != nil {
		if IsStringNumeric(val.(string)) {
			floatValue, err := StringToFloat64(val.(string))

			if err != nil {
				return 0, err
			}

			result := value + int64(floatValue)

			err = as.Put(key, result, 0)

			return result, err
		}

	}

	err := as.Put(key, value, 0)

	return value, err
}

// Decrement an integer counter by a given value
func (as *ArrayStore) Decrement(key string, value int64) (int64, error) {
	return as.Increment(key, -value)
}

// Put a value in the given store for a predetermined amount of time in mins.
func (as *ArrayStore) Put(key string, value interface{}, minutes int) error {
	val, err := Encode(value)

	mins := strconv.Itoa(minutes)

	mins = ""

	as.Client[as.GetPrefix()+key+mins] = val

	return err
}

// Put a value in the given store until it is forgotten/evicted
func (as *ArrayStore) Forever(key string, value interface{}) error {
	return as.Put(key, value, 0)
}

// Flush the store
func (as *ArrayStore) Flush() (bool, error) {
	as.Client = make(map[string]interface{})

	return true, nil
}

// Forget a given key-value pair from the store
func (as *ArrayStore) Forget(key string) (bool, error) {
	_, ok := as.Client[as.GetPrefix()+key]

	if ok {
		delete(as.Client, as.GetPrefix()+key)

		return true, nil
	}

	return false, nil
}

// Get the cache key prefix
func (as *ArrayStore) GetPrefix() string {
	return as.Prefix
}

// Put many values in the given store until they are forgotten/evicted
func (as *ArrayStore) PutMany(values map[string]interface{}, minutes int) error {
	for key, value := range values {
		as.Put(key, value, minutes)
	}

	return nil
}

// Get many values from the store
func (as *ArrayStore) Many(keys []string) (map[string]interface{}, error) {
	items := make(map[string]interface{})

	for _, key := range keys {
		val, err := as.Get(key)

		if err != nil {
			return items, err
		}

		items[key] = val
	}

	return items, nil
}

// Get the struct representation of a value from the store
func (as *ArrayStore) GetStruct(key string, entity interface{}) (interface{}, error) {
	value := as.Client[as.GetPrefix()+key]

	return Decode(value.(string), entity)
}

// Return the TaggedCache for the given store
func (as *ArrayStore) Tags(names []string) TaggedStoreInterface {
	return &TaggedCache{
		Store: as,
		Tags: TagSet{
			Store: as,
			Names: names,
		},
	}
}
