# Student Weaviate Schema & Data Migration Tools
# ================================================
# Tools for setting up and populating the student-facing Weaviate database
# with humanOS behavioral elements for action plan generation.

# Configuration
STUDENT_WEAVIATE_URL ?= http://localhost:8088
HUMANOS_PATH ?= ../../humanOS

# Export for child processes
export STUDENT_WEAVIATE_URL
export HUMANOS_PATH

.PHONY: all setup-schema migrate-all migrate-chisg migrate-barriers migrate-patterns clean help

# Default target: show help
help:
	@echo "Student Weaviate Migration Tools"
	@echo "================================="
	@echo ""
	@echo "Setup:"
	@echo "  make setup-schema    - Create all 9 Weaviate collections"
	@echo ""
	@echo "Data Migration:"
	@echo "  make migrate-all     - Run all migrations (schema + data)"
	@echo "  make migrate-chisg   - Import ~459 CHISG skills from humanOS"
	@echo "  make migrate-barriers - Import barriers and levers"
	@echo "  make migrate-patterns - Import teaching patterns"
	@echo ""
	@echo "Utilities:"
	@echo "  make clean           - Remove built binaries"
	@echo "  make verify          - Check Weaviate connection and schema"
	@echo ""
	@echo "Configuration:"
	@echo "  STUDENT_WEAVIATE_URL = $(STUDENT_WEAVIATE_URL)"
	@echo "  HUMANOS_PATH         = $(HUMANOS_PATH)"

# Build all tools
build-all:
	@echo "🔨 Building migration tools..."
	cd setup-student-schema && go build -o ../bin/setup-student-schema .
	cd migrate-chisg-skills && go build -o ../bin/migrate-chisg-skills .
	cd import-barriers && go build -o ../bin/import-barriers .
	cd import-teaching-patterns && go build -o ../bin/import-teaching-patterns .
	@echo "✅ All tools built"

# Setup Weaviate schema (9 collections)
setup-schema: 
	@echo "📋 Setting up Student Weaviate Schema..."
	@echo "   Target: $(STUDENT_WEAVIATE_URL)"
	cd setup-student-schema && go run .

# Migrate CHISG skills (~459 entries)
migrate-chisg:
	@echo "🧠 Migrating CHISG Skills..."
	@echo "   Source: $(HUMANOS_PATH)/data/chisg_elements.json"
	cd migrate-chisg-skills && go run . $(HUMANOS_PATH)/data/chisg_elements.json

# Migrate barriers and levers
migrate-barriers:
	@echo "🚧 Importing Barriers & Levers..."
	@echo "   Source: $(HUMANOS_PATH)/shared/schemas/barriers.json"
	cd import-barriers && go run . $(HUMANOS_PATH)/shared/schemas/barriers.json

# Migrate teaching patterns
migrate-patterns:
	@echo "📚 Importing Teaching Patterns..."
	cd import-teaching-patterns && go run .

# Run all migrations in order
migrate-all: setup-schema migrate-chisg migrate-barriers migrate-patterns
	@echo ""
	@echo "🎉 All migrations complete!"
	@echo ""
	@echo "Collections created:"
	@echo "  ✅ ETPProfile"
	@echo "  ✅ LearningIntervention"
	@echo "  ✅ ActionLog"
	@echo "  ✅ StudentProgress"
	@echo "  ✅ InterventionLibrary"
	@echo "  ✅ CHISGSkill (~459 skills)"
	@echo "  ✅ StudentBarrier"
	@echo "  ✅ InterventionLever"
	@echo "  ✅ TeachingPattern"

# Verify Weaviate connection and schema
verify:
	@echo "🔍 Verifying Weaviate..."
	@curl -s $(STUDENT_WEAVIATE_URL)/v1/schema | jq '.classes | length' | xargs -I {} echo "   Collections: {}"
	@curl -s $(STUDENT_WEAVIATE_URL)/v1/schema | jq '.classes[].class' 

# Clean built binaries
clean:
	rm -rf bin/
	@echo "🧹 Cleaned"

# Quick test - just setup schema
test-schema: setup-schema verify
	@echo "✅ Schema test passed"
