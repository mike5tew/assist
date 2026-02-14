import json
import urllib.request
import urllib.error

# Configuration
WEAVIATE_URL = "http://localhost:8081/v1/graphql"

# The colors defined in the Frontend (GraphExplorer.tsx / NavigationDrawer.tsx)
FRONTEND_DOMAINS = {
    'COGNITION',
    'PROBLEM SOLVING',
    'PEOPLE',
    'SELF',
    'COMMUNICATION',
    'LANGUAGE',
    'LITERACY',
    'STRATEGIC',
    'FOCUS & TOOLS',
    'NUMBERS',
    'REASONING',
    # Aliases
    'MONEY & NUMBERS',
    'CREATIVITY',
    # Legacy/Fallback
    'FOUNDATIONAL',
    'PHYSICAL',
    'default'
}

def get_weaviate_domains():
    # Trying 'name' instead of 'skillName'
    query = """
    {
      Get {
        CHISGElement(limit: 10000) {
          domain
          name
        }
      }
    }
    """
    
    req = urllib.request.Request(
        WEAVIATE_URL, 
        data=json.dumps({'query': query}).encode('utf-8'),
        headers={'Content-Type': 'application/json'}
    )

    try:
        with urllib.request.urlopen(req) as response:
            data = json.loads(response.read().decode())
            
        if 'errors' in data:
            # Fallback to just domain if name fails
            return get_weaviate_domains_fallback()

        elements = data.get('data', {}).get('Get', {}).get('CHISGElement', [])
        return process_elements(elements)
        
    except urllib.error.URLError as e:
        print(f"Failed to connect to Weaviate: {e}")
        return {}

def get_weaviate_domains_fallback():
    query = """
    {
      Get {
        CHISGElement(limit: 10000) {
          domain
        }
      }
    }
    """
    try:
        req = urllib.request.Request(
            WEAVIATE_URL, 
            data=json.dumps({'query': query}).encode('utf-8'),
            headers={'Content-Type': 'application/json'}
        )
        with urllib.request.urlopen(req) as response:
            data = json.loads(response.read().decode())
        elements = data.get('data', {}).get('Get', {}).get('CHISGElement', [])
        return process_elements(elements)
    except:
        return {}

def process_elements(elements):
    domain_counts = {}
    for el in elements:
        domain = el.get('domain')
        if not domain:
            domain = "MISSING_OR_EMPTY"
        
        if domain not in domain_counts:
            domain_counts[domain] = 0
        domain_counts[domain] += 1
    return domain_counts

def main():
    print("--- Domain Color Debugger ---")
    print(f"Checking Weaviate at: {WEAVIATE_URL}...\n")
    
    backend_domains = get_weaviate_domains()
    
    if not backend_domains:
        print("No data found or connection failed.")
        return

    print(f"{'DOMAIN IN DB':<30} | {'COUNT':<10} | {'STATUS IN FRONTEND':<20}")
    print("-" * 70)

    all_domains = set(backend_domains.keys())
    
    # Check for mismatches
    for domain in sorted(all_domains):
        count = backend_domains[domain]
        if domain in FRONTEND_DOMAINS:
            status = "✅ MATCH"
        elif domain == "MISSING_OR_EMPTY":
            status = "❌ NULL VALUE"
        elif domain.upper() in FRONTEND_DOMAINS:
             status = "⚠️ CASE MISMATCH"
        else:
            status = "❌ UNKNOWN KEY"
            
        print(f"{domain:<30} | {count:<10} | {status:<20}")

    print("\n--- Summary ---")
    
    unknown = [d for d in all_domains if d not in FRONTEND_DOMAINS and d not in ["MISSING_OR_EMPTY"]]
    if unknown:
        print("Found domains in DB that are missing in the Frontend map:")
        for u in unknown:
            print(f" - {u}")
    
    nulls = backend_domains.get("MISSING_OR_EMPTY", 0)
    if nulls > 0:
        print(f"\nWARNING: {nulls} items have no domain assigned!")
    else:
        print("All database items have a domain assigned.")

if __name__ == "__main__":
    main()
