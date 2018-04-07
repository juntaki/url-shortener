package urlshortener

import (
	"html/template"
	"net/url"
	"reflect"
	"testing"
)

func TestShortURL_Base64QRCode(t *testing.T) {
	baseURL, _ = url.Parse("https://s.juntaki.com")
	type fields struct {
		ID  string
		URL string
	}
	tests := []struct {
		name   string
		fields fields
		want   template.URL
	}{
		{
			name: "short URL properly bulid",
			fields: fields{
				ID:  "a",
				URL: "URL",
			},
			want: template.URL("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAQAAAAEACAMAAABrrFhUAAAABlBMVEX///8AAABVwtN+AAAB80lEQVR42uzbTW7CMBCA0fb+l+66EkLz42Cbed8ONUDyVCmxsX8kSZIkSZIkSZIkSZIkSZJ+s/1/37sPe3nky1fpbwcAAAAAAABWAlSPTJ9egLF5ngAAAAAAAEARIH3nD3xK+hkh/e0AAAAAAADAVoD4qBgAAAAAAAD4HoDmIQAAAAAAAMBBANUjA7f8B78dAAAAAAAA6AE01wg1X92/SAoAAAAAANwJ8MQwOj5UvjgAAAAAAPANANUfMgMTuemnggOGwwAAAAAAYBZA9bwCl5wez+4cBwMAAAAAgMkAgdt6/G/xEXOc/7P/CAAAAAAAYBbA0nt21WHn4wAAAAAAAJgMELjpxq818PbjZoUBAAAAAMBIgPjS1ep5NT9sy6+iAAAAAABgFsCa6eDmzo/jtswAAAAAAIAhANUx65qzbJ4EAAAAAAAA8BzAuyOr08jNKWYAAAAAAABga/FLjl8PAAAAAAAAsOc5oDqfW11iG3+a+NBwGAAAAAAATAZoHhmf1j18kRQAAAAAAJgMsGYDZHMZ0bsnDQAAAAAAAOB0gPT8cXy8DgAAAAAAAJwAUJ3rvWrHCAAAAAAAmAWwhio9m5zerwkAAAAAAACsBGju52gOqnf+RgoAAAAAACYDSJIkSZIkSZIkSZIkSZIk3dZfAAAA//+rmT5n+id12QAAAABJRU5ErkJggg=="),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ShortURL{
				ID:  tt.fields.ID,
				URL: tt.fields.URL,
			}
			if got := s.Base64QRCode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ShortURL.Base64QRCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShortURL_ShortURL(t *testing.T) {
	type fields struct {
		ID  string
		URL string
	}
	tests := []struct {
		name    string
		fields  fields
		baseURL string
		want    template.URL
	}{
		{
			name:    "short URL properly bulid",
			baseURL: "https://s.juntaki.com",
			fields: fields{
				ID:  "a",
				URL: "URL",
			},
			want: template.URL("https://s.juntaki.com/a"),
		},
		{
			name:    "short URL properly bulid",
			baseURL: "https://s.juntaki.com/",
			fields: fields{
				ID:  "a",
				URL: "URL",
			},
			want: template.URL("https://s.juntaki.com/a"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseURL, _ = url.Parse(tt.baseURL)
			s := &ShortURL{
				ID:  tt.fields.ID,
				URL: tt.fields.URL,
			}
			if got := s.ShortURL(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ShortURL.ShortURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
