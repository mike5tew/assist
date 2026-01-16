# Semantic Links UX Preferences

This document records the UX decisions and defaults for the semantic link extraction workflow.

Summary (defaults)
- Default source: use the **most-recent upload/source** from the upload page when beginning extraction.
- Source change: user may change the source via a small searchable dialog accessible from the extraction page.
- Selection model: 3-click/3-step workflow (Source term -> Target term -> Relationship selection). Quick-select relationship list is shown after selecting both terms.
- Relationship list: a short list of frequent relation types (e.g., causes, treats, is_a, part_of, requires_for, enables) with an **Add relationship** option.
- Persistence: new relations created by the user are persisted client-side (and can later be synced to server if approved).
- Preview: PDF/text preview must be visible so that terms can be selected directly from page text (client-side preview used as fallback when server-processed text is not present).
- Keyboard accessibility: text selection and relation selection should be keyboard-accessible; consider shortcuts for "select source" (S) and "select target" (T).

Detailed notes
- Defaulting source from upload: When the upload flow completes or the user opens the extractor from an upload confirmation, the extractor should prefills the "source document" field with the uploaded document and allow override via dialog. This reduces friction.
- Selection workflow: Selection should be as fast as possible: selecting a phrase in the preview should immediately populate the Source or Target field and present visually which one is selected. If both are selected, present relationship quick-list and "Create Link" CTA.
- Relationship quick-list: show the top 6 relations by default with an "Add more" button that opens a small modal for creating custom relation types. Allow quick persistence of newly created types.
- Fallbacks: If the backend has processed text for a document, prefer loading that (more accurate tokenization and positions). If not, use `pdfjs-dist` client-side text extraction for immediate selection.
- Accessibility: ARIA labels for preview container, explicit focus and keyboard handlers for selecting source/target, and descriptive tooltips for relationship buttons.

Acceptance criteria
- The extractor pre-selects the source when opened from an upload.
- Selecting text in the preview sets Source and then Target in order.
- Relationship quick-list is fast to use and new relations can be added via modal.
- The flow is documented here and referenced in the PR description when changes are made.
