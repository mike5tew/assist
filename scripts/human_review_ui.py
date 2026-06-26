import streamlit as st
import json
import os
import argparse
from pathlib import Path

def load_data(file_path):
    if not os.path.exists(file_path):
        return []
    with open(file_path, "r", encoding="utf-8") as f:
        return json.load(f)

def save_data(data, file_path):
    with open(file_path, "w", encoding="utf-8") as f:
        json.dump(data, f, indent=2)

def main():
    st.set_page_config(page_title="CHISG Human Review", layout="wide")
    st.title("CHISG Link Review Interface")

    # Sidebar for file selection
    st.sidebar.header("Configuration")
    data_dir = st.sidebar.text_input("Data Directory", value=".")
    
    json_files = [f for f in os.listdir(data_dir) if f.endswith("_extracted.json")]
    if not json_files:
        st.warning(f"No extracted JSON files found in {os.path.abspath(data_dir)}")
        return

    selected_file = st.sidebar.selectbox("Select File to Review", json_files)
    file_path = os.path.join(data_dir, selected_file)
    approved_file_path = os.path.join(data_dir, selected_file.replace("_extracted.json", "_approved.json"))

    if "current_file" not in st.session_state or st.session_state.current_file != file_path:
        st.session_state.current_file = file_path
        st.session_state.links = load_data(file_path)
        st.session_state.approved_links = load_data(approved_file_path)

    links = st.session_state.links
    if not links:
        st.info("No links found in this file.")
        return

    # Metrics
    approved_ids = {f"{l.get('source_term')}::{l.get('target_term')}" for l in st.session_state.approved_links}
    st.sidebar.metric("Total Links", len(links))
    st.sidebar.metric("Approved", len(approved_ids))
    st.sidebar.metric("Pending", len(links) - len(approved_ids))

    st.write(f"### Reviewing: {selected_file}")

    # Display Links
    for i, link in enumerate(links):
        link_id = f"{link.get('source_term')}::{link.get('target_term')}"
        is_approved = link_id in approved_ids

        with st.expander(f"{'✅' if is_approved else '⏳'} {link.get('source_term')} -> {link.get('relation_type')} -> {link.get('target_term')}", expanded=not is_approved):
            cols = st.columns([1, 1])
            
            with cols[0]:
                st.markdown("**Source Term:**")
                st.text_input("Source", value=link.get('source_term'), key=f"src_{i}")
                
                st.markdown("**Relation:**")
                st.selectbox("Relation", options=[link.get('relation_type'), "causes", "enables", "inhibits", "is composed of", "is a type of", "has property", "has value", "is used for", "is found in", "is an example of", "determines", "represents", "contradicts", "correlates with", "supports", "refutes", "regulates", "is a mechanism of", "is a marker for", "interacts with"], key=f"rel_{i}")
                
                st.markdown("**Target Term:**")
                st.text_input("Target", value=link.get('target_term'), key=f"tgt_{i}")
                
                st.markdown("**Thread ID:**")
                st.text_input("Thread", value=link.get('thread_id'), key=f"thr_{i}")

            with cols[1]:
                st.markdown("**Context:**")
                st.text_area("Context", value=link.get('context'), key=f"ctx_{i}", height=100)
                
                st.markdown("**Explanation / Rationale:**")
                st.text_area("Rationale", value=f"Explanation: {link.get('explanation')}\n\nHierarchy: {link.get('hierarchy_rationale')}", key=f"rat_{i}", height=100)
                
                m1, m2, m3, m4 = st.columns(4)
                m1.metric("Src Gen", link.get('source_term_generality'))
                m2.metric("Tgt Gen", link.get('target_term_generality'))
                m3.metric("Sem Dist", link.get('semantic_distance'))
                m4.metric("Strength", link.get('relationship_strength'))

            # Approval action
            if st.button("Approve & Save" if not is_approved else "Update Approved", key=f"btn_{i}"):
                # Update link with local edits
                link['source_term'] = st.session_state[f"src_{i}"]
                link['relation_type'] = st.session_state[f"rel_{i}"]
                link['target_term'] = st.session_state[f"tgt_{i}"]
                link['thread_id'] = st.session_state[f"thr_{i}"]
                link['context'] = st.session_state[f"ctx_{i}"]
                
                # Add to approved list avoiding duplicates
                current_approved = [l for l in st.session_state.approved_links if f"{l.get('source_term')}::{l.get('target_term')}" != link_id]
                current_approved.append(link)
                st.session_state.approved_links = current_approved
                save_data(current_approved, approved_file_path)
                st.rerun()

if __name__ == "__main__":
    main()
