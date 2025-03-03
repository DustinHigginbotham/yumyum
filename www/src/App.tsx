import { useState } from 'react';
import './App.css';

import { useRecipeGenerator } from './hooks/useRecipeGenerator';

function App() {
  const [ name, setName ] = useState('');
  const [ ingredients, setIngredients ] = useState('');
  
  const {
    error,
    story,
    isLoading,
    isWriting,
    generateRecipe,
    copyToClipboard,
  } = useRecipeGenerator();

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    generateRecipe(name, ingredients);
  }

  return (
    <>


      <h1>Yum Yum!</h1>
      <p>Crafting the perfect recipe backstory can be time-consuming, and let's be honest, not everyone has a heartwarming tale about artisanal rosemary or a life-changing bowl of soup. But don't worry, <em>we've got you covered</em>!</p>

      <form onSubmit={handleSubmit}>
        <div className="yummy-input-row name">
          <label htmlFor="name">Give your recipe a name!</label>
          <input
            type="text"
            name="name"
            id="name"
            placeholder="Grandma's Famous Lasagna"
            value={name}
            onInput={e => setName((e.target as HTMLInputElement).value)} />
        </div>

        <div className="yummy-input-row ingredients">
          <label htmlFor="ingredients">List the ingredients, one per line</label>
          <textarea
            name="ingredients"
            id="ingredients"
            placeholder="1 cup flour"
            value={ingredients}
            onInput={e => setIngredients((e.target as HTMLTextAreaElement).value)} />
        </div>

        <div className="yummy-input-row submit">
          <button type="submit" disabled={isLoading}>Generate</button>
          {error && <div className="yummy-error">There was an error generating your recipe backstory!</div>}
        </div>
      </form>

      {isLoading && <div className="yummy-loading">Generating...</div>}
      {story && (
        <>
          <div className={`yummy-story-time ${isWriting ? 'writing' : 'done'}`} dangerouslySetInnerHTML={{ __html: story }} />
          <button disabled={isWriting} onClick={copyToClipboard}>Copy</button>
        </>
      )}
    </>
  )
}

export default App
