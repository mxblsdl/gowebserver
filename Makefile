.PHONY: dev templ-watch stop

# Start both air and templ watch in parallel
dev: 
	air 
	templ-watch

# Run air for hot-reloading
air:
	air

# Run templ in watch mode
templ-watch:
	templ generate --watch

# Stop all development processes
stop:
	@echo "Stopping development servers..."
	@pkill -f "air" || true
	@pkill -f "templ" || true
	@echo "Development servers stopped"

# Default target runs both processes
all: dev