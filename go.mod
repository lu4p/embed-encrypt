module github.com/abakum/embed-encrypt

go 1.21.4

replace internal/tool => ./example/internal/tool

replace public/tool => ./example/public/tool

require (
	internal/tool v0.0.0-00010101000000-000000000000
	public/tool v0.0.0-00010101000000-000000000000
)
