// VoiceQueryComponent.tsx
import React, { useState } from 'react';

const VoiceQueryComponent: React.FC = () => {
    const [isListening, setIsListening] = useState(false);
    const [transcript, setTranscript] = useState('');
    const [answer, setAnswer] = useState('');

    // Initialize Speech Recognition
    const SpeechRecognition = (window as any).SpeechRecognition || (window as any).webkitSpeechRecognition;
    const recognition = new SpeechRecognition();
    recognition.continuous = false;
    recognition.lang = 'en-US';
    recognition.interimResults = false; // We only want final results

    recognition.onresult = (event: any) => {
        const userSpeech = event.results[0][0].transcript;
        setTranscript(userSpeech);
        setIsListening(false);
        
        // Process the query
        handleVoiceQuery(userSpeech);
    };

    recognition.onerror = (event: any) => {
        console.error('Speech recognition error', event.error);
        setIsListening(false);
    };

    const startListening = () => {
        setTranscript('');
        setAnswer('');
        setIsListening(true);
        recognition.start();
    };

    const handleVoiceQuery = async (userTranscript: string) => {
        try {
            // Step 1: Send to Adobe API for query refinement
            const refinedQuery = await callAdobeLLM(userTranscript);
            
            // Step 2: Send refined query to your HSG Backend
            const hsgResponse = await queryHSGBackend(refinedQuery);
            
            // Step 3: Display and speak the answer
            setAnswer(hsgResponse.answer);
            speakText(hsgResponse.answer);
            
        } catch (error) {
            console.error('Query processing failed:', error);
            setAnswer('Sorry, I encountered an error processing your question.');
        }
    };

    const callAdobeLLM = async (userTranscript: string): Promise<string> => {
        // Use the precise prompt from Section 1 above
        const adobePrompt = `ROLE: You are a specialized query pre-processor... USER'S SPOKEN QUESTION TO PROCESS: "${userTranscript}"`;
        
        const response = await fetch('/api/adobe-refine-query', { // Your proxy endpoint to Adobe API
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ prompt: adobePrompt })
        });
        
        const data = await response.json();
        return data.refinedQuery; // Expecting just the string back
    };

    const queryHSGBackend = async (refinedQuery: string) => {
        const response = await fetch('/api/hsg-query', { // Your existing HSG endpoint
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ query: refinedQuery })
        });
        return await response.json();
    };

    const speakText = (text: string) => {
        const utterance = new SpeechSynthesisUtterance(text);
        window.speechSynthesis.speak(utterance);
    };

    return (
        <div>
            <button onClick={startListening} disabled={isListening}>
                {isListening ? 'Listening...' : 'Start Voice Query'}
            </button>
            {transcript && <p><strong>You asked:</strong> {transcript}</p>}
            {answer && <p><strong>Answer:</strong> {answer}</p>}
        </div>
    );
};

export default VoiceQueryComponent;