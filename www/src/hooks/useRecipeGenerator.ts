import { useState, useEffect, useCallback, useRef } from 'react';
import { marked } from 'marked';

marked.setOptions({
    gfm: true,
    breaks: true,
});

const API_URL = import.meta.env.VITE_API_URL;

export function useRecipeGenerator() {
    const [ story, setStory ] = useState('');
    const [ isLoading, setIsLoading ] = useState(false);
    const [ isWriting, setIsWriting ] = useState(false);
    const [ eventSource, setEventSource ] = useState<EventSource | null>(null);
    const [ error, setError ] = useState<any>(null);

    const accumulatedStoryRef = useRef('');

    let accumulatedStory = '';

    useEffect(() => {
        if (!eventSource) return;

        eventSource.addEventListener('message', async (ev: MessageEvent) => {

            const { data } = ev;

            if (isLoading) setIsLoading(false);
            if (!isWriting) setIsWriting(true);

            setError(null);

            if (data === '[DONE]') {
                setIsWriting(false);
                eventSource.close();
                setEventSource(null);
                return;
            }

            accumulatedStoryRef.current += data.replaceAll('[NEWLINE]', '\n');
            setStory(await marked.parse(accumulatedStoryRef.current));
        });

        eventSource.addEventListener('error', (error: Event) => {
            setError(error);
            console.error('EventSource error:', error);
            setIsWriting(false);
            setIsLoading(false);
            eventSource.close();
            setEventSource(null);
        });

        return () => {
            if (eventSource) {
                eventSource.close();
            }
        };
    }, [eventSource]);

    const generateRecipe = useCallback(async (name: string, ingredients: string) => {
        setIsLoading(true);
        setStory('');
        accumulatedStoryRef.current = '';

        const formattedIngredients = ingredients.split('\n').join(',');
        const newEventSource = new EventSource(`${API_URL}/generate?name=${name}&ingredients=${formattedIngredients}`);

        setEventSource(newEventSource);
    }, []);

    const copyToClipboard = useCallback(() => {
        navigator.clipboard.writeText(accumulatedStoryRef.current);
    }, [accumulatedStory]);

    return {
        story,
        isWriting,
        isLoading,
        error,
        generateRecipe,
        copyToClipboard,
    };
}