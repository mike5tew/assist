// MongoDB diagnostic script - run with: mongo diagnostic.js

// Function to print collection info
function printCollectionInfo(dbName, collName) {
  print("\n---- Collection: " + collName + " ----");
  var count = db.getSiblingDB(dbName).getCollection(collName).countDocuments({});
  print("Document count: " + count);
  
  if (count > 0) {
    print("Sample document:");
    var sample = db.getSiblingDB(dbName).getCollection(collName).findOne();
    printjson(sample);
    
    // Check for text fields that might contain "xla" or "agammaglobulinemia"
    print("\nSearching for 'xla' or 'agammaglobulinemia':");
    var xlaCount = db.getSiblingDB(dbName).getCollection(collName).countDocuments({
      $or: [
        { content: { $regex: "xla", $options: "i" } },
        { content: { $regex: "agammaglobulinemia", $options: "i" } },
        { chapter_title: { $regex: "xla", $options: "i" } },
        { chapter_title: { $regex: "agammaglobulinemia", $options: "i" } },
        { term: { $regex: "xla", $options: "i" } },
        { term: { $regex: "agammaglobulinemia", $options: "i" } }
      ]
    });
    print("Documents mentioning XLA or agammaglobulinemia: " + xlaCount);
    
    // Show fields present in collection
    print("\nFields present in collection:");
    var fields = {};
    db.getSiblingDB(dbName).getCollection(collName).find().limit(10).forEach(function(doc) {
      Object.keys(doc).forEach(function(key) {
        fields[key] = true;
      });
    });
    printjson(Object.keys(fields));
  }
}

// List all databases
print("===== AVAILABLE DATABASES =====");
db.adminCommand('listDatabases').databases.forEach(function(d) {
  print(d.name);
});

// Check esp_organizer database
print("\n===== ESP_ORGANIZER DATABASE =====");
var dbName = "esp_organizer";
var collNames = ["immunology_terms", "immunology_chapters", "immunology_case_studies", "skills"];

collNames.forEach(function(collName) {
  printCollectionInfo(dbName, collName);
});

// Check esp_project database
print("\n===== ESP_PROJECT DATABASE =====");
dbName = "esp_project";
collNames.forEach(function(collName) {
  printCollectionInfo(dbName, collName);
});

// Check skills_dev database
print("\n===== SKILLS_DEV DATABASE =====");
dbName = "skills_dev";
collNames.forEach(function(collName) {
  printCollectionInfo(dbName, collName);
});
