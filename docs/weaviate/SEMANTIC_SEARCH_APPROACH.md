# Semantic Search (Weaviate) — Current Approach + Long-Term Notes

This doc records how semantic search is currently implemented for CHISG/skills queries in the `assist` stack, why it required a vector backfill, and what needs to be true long-term for this to stay reliable.

## Summary

- The API endpoint `POST /api/skills/semantic-query` is implemented in `assist/esp-organizer`.
- It queries MongoDB for exact/keyword matches, and Weaviate for semantic matches.
- The active Weaviate instance for CHISG data is the **humanOS Weaviate** published on `http://localhost:8081` (from the host), not the legacy compose Weaviate on `8088`.
- The CHISG classes in that Weaviate are configured with `vectorizer: none`, which means:
  - `nearText` may be unavailable for those classes.
  - `nearVector` only works if objects already have stored vectors.
- Because vectors were missing, semantic search initially returned empty results until vectors were backfilled.

## Why `nearVector` (not `nearText`)

In the humanOS Weaviate schema (the one that actually contains `CHISGElement` objects), the class is configured with `vectorizer: none`. In this setup:

- We cannot rely on server-side text vectorization.
- We generate embeddings in the API and query via `nearVector`.

## Why semantic search was empty

Even after switching to `nearVector`, results were empty because the underlying `CHISGElement` objects had **no stored vectors**.

With `vectorizer: none`, Weaviate does not generate/store vectors automatically. If the data was imported without explicit vectors, `nearVector` has nothing to compare against.

## Operational fix: vector backfill

A dev-only endpoint was added to backfill vectors into Weaviate for existing objects:

- `POST /api/diagnostics/weaviate/backfill-vectors`
- Guarded by `ENABLE_WEAVIATE_VECTOR_BACKFILL=true`

Backfill behavior:

- Fetches objects for a class (initially `CHISGElement`).
- Skips objects that already have vectors (idempotent).
- Generates an embedding for each object’s derived text.
- Writes vectors back into Weaviate.

After this ran, semantic search began returning results.

## Add-doc endpoint (optional tooling)

A second dev-only endpoint exists to insert a single `Documentation` object into Weaviate with an explicit vector:

- `POST /api/diagnostics/weaviate/add-doc`
- Guarded by `ENABLE_WEAVIATE_VECTOR_BACKFILL=true`

This is intended for writing internal “what changed” notes into the `Documentation` class in the same vectorless schema.

## Configuration: which Weaviate are we using?

In `assist/docker-compose.yml` the `assist-api` container is currently configured to use the humanOS Weaviate:

- `WEAVIATE_URL: http://host.docker.internal:8081`

This is a workaround for the fact that the humanOS Weaviate runs on a different compose/network.

## Long-term sufficiency checklist

This approach is workable, but only if we treat vector storage as part of ingestion and ongoing maintenance.

Long-term, one of these needs to be true:

1) **Ingestion always writes vectors**
- Any new `CHISGElement`/`Documentation` objects must be created with vectors already present.
- Re-embedding/re-indexing becomes an explicit operation when the embedding model changes.

2) **Enable a server-side vectorizer module and use `nearText`**
- Configure Weaviate classes to use a vectorizer module.
- Reindex existing objects.
- This reduces custom embedding logic in the API, but increases Weaviate module coupling.

3) **Consolidate Weaviate topology**
- Today, the meaningful data lives on 8081 (humanOS), while the `assist` compose still ships a legacy Weaviate on 8088.
- Reducing this to a single clearly-owned instance (or a clearly-owned pair: one for skills, one for docs) will avoid future “empty instance” regressions.

## Risks to document/monitor

- Embedding model changes: vector dimension and semantic space can change; mixing embeddings will degrade results.
- Backfill needs: if ingestion forgets to write vectors, semantic search silently degrades.
- Network topology: `host.docker.internal` is a convenience, not a robust cross-compose service discovery mechanism.

## Operational / Deployment notes (contact endpoint & secrets)

- A protected contact endpoint is available at `POST /api/contact`. It requires a reCAPTCHA token when `RECAPTCHA_SECRET` is set and will deliver messages via SMTP when SMTP secrets are configured.
- SMTP credentials are now provided via Docker secrets (`smtp_user` / `smtp_pass`) in production; the compose file expects `external: true` secrets. For local development `./secrets/*` files may be used (gitignored), but **remove** them after creating swarm/production secrets.
- Build-time frontend reCAPTCHA key: `REACT_APP_RECAPTCHA_SITE_KEY` (inject at build); server-side verification uses `RECAPTCHA_SECRET`.
- Recommended: store all production secrets in a managed secret store (Docker Swarm secrets, Vault, or cloud provider secrets) and remove plaintext from all repo/config files.
