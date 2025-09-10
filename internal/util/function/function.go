package function

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"math"
	"reflect"

	"github.com/google/uuid"
)

func GenerateUUID() string {
	return uuid.New().String()
}

func GenerateRandomString(l int) string {
	buff := make([]byte, int(math.Ceil(float64(l)/float64(1.33333333333))))
	rand.Read(buff)
	str := base64.RawURLEncoding.EncodeToString(buff)
	return str[:l] // get desired length
}

func GetTypeCount(i any) int {
	switch reflect.ValueOf(i).Kind() {
	case reflect.Map:
		return reflect.ValueOf(i).Len()
	case reflect.Array:
		return reflect.ValueOf(i).Len()
	case reflect.Slice:
		return reflect.ValueOf(i).Len()
	default:
		return 1
	}
}

func ValidJson(jsonValue json.RawMessage) bool {
	bValue, err := jsonValue.MarshalJSON()
	if err != nil {
		return false
	}
	check := make(map[string]any, 0)
	if errCheck := json.Unmarshal(bValue, &check); errCheck != nil {
		check := make([]any, 0)
		if errCheck := json.Unmarshal(bValue, &check); errCheck != nil {
			// if it is not a map or a slice, then it is not valid JSON
			return false
		}
	}
	return true
}

// for each element in 'compare', if NOT in 'src', it will be added to the resulting 'diff'
// if you want to add from an existing, compared to another list, src => existing; compare => new list
// if you want to delete from to an existing, compared to another list, src => new list; compare => existing
func ArrayDiff(src, compare []string) (diff []string) {
	m := make(map[string]struct{})
	for _, i := range src {
		m[i] = struct{}{}
	}
	for _, i := range compare {
		if _, ok := m[i]; !ok {
			diff = append(diff, i)
		}
	}
	return
}
