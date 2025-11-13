#!/bin/bash
# scripts/seed-medical-data.sh
# This script seeds minimal test data into the medical collections

set -e
echo "🧪 Seeding test medical data..."

# Connect to MongoDB container
docker exec -it esp_mongodb mongosh \
  --username "$MONGO_INITDB_ROOT_USERNAME" \
  --password "$MONGO_INITDB_ROOT_PASSWORD" \
  --authenticationDatabase admin \
  --eval '
    use esp_organizer;
    
    // Insert sample case studies
    db.case_studies.deleteMany({});
    db.case_studies.insertMany([
      {
        "case_number": "Case 1",
        "title": "X-Linked Agammaglobulinemia",
        "content": "A 2-year-old boy presents with recurrent sinopulmonary infections and absence of B cells.",
        "clinical_findings": ["Recurrent bacterial infections", "Absence of B cells", "Low immunoglobulins"]
      },
      {
        "case_number": "Case 2",
        "title": "Severe Combined Immunodeficiency",
        "content": "A 6-month-old presents with failure to thrive and opportunistic infections.",
        "clinical_findings": ["Opportunistic infections", "Lymphopenia", "Failure to thrive"]
      }
    ]);
    
    // Insert sample medical terms
    db.medical_terms.deleteMany({});
    db.medical_terms.insertMany([
      {
        "term": "Agammaglobulinemia",
        "category": "Immunodeficiency",
        "context": ["Absence of gamma globulins (antibodies) in the blood"]
      },
      {
        "term": "B-cell",
        "category": "Cellular",
        "context": ["White blood cells that produce antibodies"]
      },
      {
        "term": "BTK gene",
        "category": "Genetic",
        "context": ["Gene mutated in X-linked agammaglobulinemia"]
      }
    ]);
    
    // Verify data was inserted
    print("✅ Case studies inserted: " + db.case_studies.countDocuments());
    print("✅ Medical terms inserted: " + db.medical_terms.countDocuments());
  '

echo "✅ Medical data seeding complete"
