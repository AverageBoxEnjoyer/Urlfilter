<script>
  import { onMount } from 'svelte';
  import { SearchCredentials, GetVersion, FetchApiMessage, CheckCombosFolder } from '../wailsjs/go/main/App';
  
  let keyword = "";
  let results = {
    count: 0,
    duration: "",
    filename: ""
  };
  let isLoading = false;
  let version = "";
  let errorMessage = "";
  let searchHistory = [];
  let showHistory = false;
  
  // Mensaje para la barra de anuncios
  let announcementMessage = "No olvides unirte al canal y chat. | https://t.me/ConfigsINCs & https://t.me/ConfigsIncsChat";
  let showAnnouncement = true;

// Añade estas variables junto a las otras variables existentes
// Reemplaza las variables y funciones relacionadas con la API
// Añade estas variables en la sección de script
let apiMessage = "";
let apiLoading = true;

function dismissError() {
  errorMessage = '';
}

// Función simple para obtener datos de la API
async function fetchApiData() {
    apiLoading = true;
    
    try {
      const jsonString = await FetchApiMessage();
      console.log("Respuesta de la API:", jsonString);
      
      // Parsear el JSON
      const data = JSON.parse(jsonString);
      const mensaje = data.message; 
      
      // Actualizar el mensaje
      apiMessage = mensaje;
      showAnnouncement = true;
    } catch (error) {
      console.error("Error al obtener datos de la API:", error);
      apiMessage = announcementMessage;
    } finally {
      apiLoading = false;
    }
  }
  
  onMount(async () => {
    try {
      version = await GetVersion();
    } catch (error) {
      console.error("Error al obtener versión:", error);
    }
    
    fetchApiData();
  });

  async function search() {
    // Limpiar errores anteriores al iniciar nueva búsqueda
    errorMessage = "";
    try {
    await CheckCombosFolder();
  } catch (err) {
    errorMessage = typeof err === 'object' ? err.message : String(err);
    console.error('Raw error:', err);
    return;
  }
    
    if (!keyword.trim()) {
        errorMessage = "Por favor ingresa una palabra clave";
        return;
    }

    try {
        await CheckCombosFolder();
    } catch (err) {
        errorMessage = err.message;
        console.log("Error CheckCombosFolder:", err); // Log detallado
        return;
    }
    
    isLoading = true;
    
    try {
        const response = await SearchCredentials(keyword);
        results = {
            count: response.count || 0,
            duration: response.duration || "0s",
            filename: response.filename || ""
        };
        
        // Actualizar historial
        searchHistory = [
            { keyword, timestamp: new Date().toLocaleTimeString(), count: results.count },
            ...searchHistory.filter(item => item.keyword !== keyword)
        ].slice(0, 5);
        
    } catch (err) {
        errorMessage = err.message; // Corregido: usar .message en lugar de err()
        console.error("Error SearchCredentials:", err);
    } finally {
        isLoading = false;
    }
}

  function handleKeydown(event) {
    if (event.key === 'Enter' && !isLoading) {
      search();
    }
  }

  function selectFromHistory(item) {
    keyword = item.keyword;
    showHistory = false;
  }
  
  function dismissAnnouncement() {
    showAnnouncement = false;
  }
</script>

<main>
  <div class="dashboard">
    <header>
      <div class="header-content">
        <h1>Dashboard de Búsqueda</h1>
        <div class="version-badge">v{version}</div>
      </div>
    </header>
    
    {#if showAnnouncement}
  <div class="announcement-bar">
    <div class="announcement-content">
      <svg class="announcement-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="12" cy="12" r="10"></circle>
        <line x1="12" y1="8" x2="12" y2="12"></line>
        <line x1="12" y1="16" x2="12.01" y2="16"></line>
      </svg>
      <p>{apiMessage}</p>
    </div>
    <button class="announcement-close" on:click={dismissAnnouncement} aria-label="Cerrar anuncio">
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <line x1="18" y1="6" x2="6" y2="18"></line>
        <line x1="6" y1="6" x2="18" y2="18"></line>
      </svg>
    </button>
  </div>
{/if}

    <div class="search-section">
      <div class="search-container">
        <div class="search-input-wrapper">
          <svg class="search-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="11" cy="11" r="8"></circle>
            <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
          </svg>
          <input 
            type="text" 
            bind:value={keyword}
            placeholder="Buscar palabra clave en URLs..."
            class="search-input"
            disabled={isLoading}
            on:keydown={handleKeydown}
            on:focus={() => showHistory = true}
            on:blur={() => setTimeout(() => showHistory = false, 200)}
          />
          {#if showHistory && searchHistory.length > 0}
            <div class="search-history">
              {#each searchHistory as item}
                <button class="history-item" on:click={() => selectFromHistory(item)}>
                  <span>{item.keyword}</span>
                  <span class="history-count">{item.count} resultados</span>
                </button>
              {/each}
            </div>
          {/if}
        </div>
        <button 
          on:click={search} 
          class="search-button"
          disabled={isLoading}
          aria-label="Buscar"
        >
          {#if isLoading}
            <div class="loader"></div>
          {:else}
            <span>Buscar</span>
          {/if}
        </button>
      </div>

      {#if errorMessage}
  <div class="error-bar">
    <div class="error-content">
      <svg class="error-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path>
  <line x1="12" y1="9" x2="12" y2="13"></line>
  <line x1="12" y1="17" x2="12.01" y2="17"></line>
</svg>
      <p>{errorMessage}</p>
    </div>
    <button class="error-close" on:click={dismissError} aria-label="Cerrar anuncio">
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <line x1="18" y1="6" x2="6" y2="18"></line>
        <line x1="6" y1="6" x2="18" y2="18"></line>
      </svg>
    </button>
  </div>
{/if}


    </div>
    <div class="results-section">
      {#if isLoading}
        <div class="loading-state">
          <div class="loader-large"></div>
          <p>Buscando resultados...</p>
        </div>
      {:else if results.count > 0}
        <h2 class="results-title">Resultados de la búsqueda</h2>
        <div class="results-grid">
          <div class="result-card">
            <div class="card-icon">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <line x1="4" y1="9" x2="20" y2="9"></line>
                <line x1="4" y1="15" x2="20" y2="15"></line>
                <line x1="10" y1="3" x2="8" y2="21"></line>
                <line x1="16" y1="3" x2="14" y2="21"></line>
              </svg>
            </div>
            <div class="card-content">
              <h3>Coincidencias</h3>
              <p class="card-value">{results.count}</p>
            </div>
          </div>
          
          <div class="result-card">
            <div class="card-icon">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="10"></circle>
                <polyline points="12 6 12 12 16 14"></polyline>
              </svg>
            </div>
            <div class="card-content">
              <h3>Tiempo de ejecución</h3>
              <p class="card-value">{results.duration}</p>
            </div>
          </div>
          
          <div class="result-card">
            <div class="card-icon">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                <polyline points="14 2 14 8 20 8"></polyline>
                <line x1="16" y1="13" x2="8" y2="13"></line>
                <line x1="16" y1="17" x2="8" y2="17"></line>
                <polyline points="10 9 9 9 8 9"></polyline>
              </svg>
            </div>
            <div class="card-content">
              <h3>Archivo generado</h3>
              <p class="card-value filename">{results.filename}</p>
            </div>
          </div>
        </div>
      {:else if !isLoading && keyword}
        <div class="no-results">
          <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <ellipse cx="12" cy="5" rx="9" ry="3"></ellipse>
            <path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"></path>
            <path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"></path>
          </svg>
          <p>No se encontraron resultados para "{keyword}"</p>
        </div>
      {:else}
        <div class="empty-state">
          <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="11" cy="11" r="8"></circle>
            <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
          </svg>
          <p>Ingresa una palabra clave para comenzar la búsqueda</p>
        </div>
      {/if}
    </div>
  </div>
</main>

<style>
  :global(body) {
    font-family: 'Inter', 'Segoe UI', sans-serif;
    background-color: #121826;
    margin: 0;
    padding: 0;
    color: #e2e8f0;
    height: 100vh;
    overflow-x: hidden;
  }

  main {
    min-height: 100vh;
    display: flex;
    flex-direction: column;
  }

  .dashboard {
    max-width: 1200px;
    width: 100%;
    margin: 0 auto;
    padding: 2rem;
    box-sizing: border-box;
  }

  header {
    margin-bottom: 1.5rem;
    border-bottom: 1px solid #2d3748;
    padding-bottom: 1rem;
  }

  .header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  h1 {
    margin: 0;
    font-size: 1.75rem;
    font-weight: 600;
    color: #f7fafc;
  }

  .version-badge {
    background-color: #2d3748;
    color: #a0aec0;
    padding: 0.25rem 0.75rem;
    border-radius: 1rem;
    font-size: 0.75rem;
    font-weight: 500;
  }
  
  /* Estilos para la barra de anuncios */
  .announcement-bar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    background-color: rgba(79, 70, 229, 0.1);
    border-left: 3px solid #4f46e5;
    border-radius: 0.375rem;
    padding: 0.75rem 1rem;
    margin-bottom: 1.5rem;
    transition: all 0.3s ease;
  }
  
  .announcement-content {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }
  
  .announcement-icon {
    color: #4f46e5;
    flex-shrink: 0;
  }
  
  .announcement-content p {
    margin: 0;
    font-size: 0.875rem;
    color: #e2e8f0;
  }
  
  .announcement-close {
    background: none;
    border: none;
    color: #a0aec0;
    cursor: pointer;
    padding: 0.25rem;
    border-radius: 0.25rem;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s ease;
  }
  
  .announcement-close:hover {
    color: #e2e8f0;
    background-color: rgba(255, 255, 255, 0.1);
  }

  .search-section {
    margin-bottom: 2rem;
  }

  .search-container {
    display: flex;
    gap: 0.75rem;
    margin-bottom: 1rem;
  }

  .search-input-wrapper {
    position: relative;
    flex-grow: 1;
  }

  .search-icon {
    position: absolute;
    left: 1rem;
    top: 50%;
    transform: translateY(-50%);
    color: #718096;
  }

  .search-input {
    width: 93%;
    padding: 0.75rem 1rem 0.75rem 2.5rem;
    background-color: #1a202c;
    color: #e2e8f0;
    border: 1px solid #2d3748;
    border-radius: 0.375rem;
    font-size: 0.875rem;
    transition: all 0.2s ease;
  }

  .search-input:focus {
    outline: none;
    border-color: #4299e1;
    box-shadow: 0 0 0 2px rgba(66, 153, 225, 0.2);
  }

  .search-input::placeholder {
    color: #718096;
  }

  .search-input:disabled {
    opacity: 0.7;
    cursor: not-allowed;
  }

  .search-history {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    background-color: #1a202c;
    border: 1px solid #2d3748;
    border-top: none;
    border-radius: 0 0 0.375rem 0.375rem;
    z-index: 10;
    max-height: 200px;
    overflow-y: auto;
  }

  .history-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 100%;
    padding: 0.75rem 1rem;
    text-align: left;
    background: none;
    border: none;
    color: #e2e8f0;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .history-item:hover {
    background-color: #2d3748;
  }

  .history-count {
    font-size: 0.75rem;
    color: #a0aec0;
  }

.search-button {
    background-color: #4f46e5; /* Azul eléctrico suave */
    color: white;
    border: none;
    border-radius: 0.375rem;
    padding: 0.5rem 1.5rem;
    font-size: 0.875rem;
    font-weight: 500;
    min-width: 100px;
    transition: all 0.2s ease;
    cursor: pointer;
}

.search-button:hover {
    background-color: #4338ca; /* Azul más profundo */
    transform: translateY(-1px);
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
}

.search-button:focus {
    outline: none;
    box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.3); /* Anillo de enfoque sutil */
}


  .results-section {
    background-color: #1a202c;
    border-radius: 0.5rem;
    padding: 1.5rem;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  }

  .results-title {
    margin-top: 0;
    margin-bottom: 1.5rem;
    font-size: 1.25rem;
    font-weight: 600;
    color: #f7fafc;
  }

  .results-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 1.5rem;
  }

  .result-card {
    background-color: #2d3748;
    border-radius: 0.5rem;
    padding: 1.5rem;
    display: flex;
    align-items: flex-start;
    gap: 1rem;
    transition: transform 0.2s, box-shadow 0.2s;
  }

  .result-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 12px rgba(0, 0, 0, 0.15);
  }

  .card-icon {
    background-color: #4299e1;
    color: white;
    width: 48px;
    height: 48px;
    border-radius: 0.375rem;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }

  .result-card:nth-child(2) .card-icon {
    background-color: #48bb78;
  }

  .result-card:nth-child(3) .card-icon {
    background-color: #ed8936;
  }

  .card-content {
    flex-grow: 1;
  }

  .card-content h3 {
    margin: 0 0 0.5rem 0;
    font-size: 0.875rem;
    font-weight: 500;
    color: #a0aec0;
  }

  .card-value {
    margin: 0;
    font-size: 1.25rem;
    font-weight: 600;
    color: #f7fafc;
  }

  .filename {
    word-break: break-all;
    font-size: 1rem;
  }

  .no-results, .empty-state, .loading-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 3rem 1rem;
    text-align: center;
    color: #a0aec0;
  }

  .no-results p, .empty-state p, .loading-state p {
    margin-top: 1rem;
    font-size: 1rem;
  }

  .loader {
    width: 18px;
    height: 18px;
    border: 2px solid rgba(255, 255, 255, 0.3);
    border-radius: 50%;
    border-top-color: white;
    animation: spin 1s linear infinite;
  }

  .loader-large {
    width: 40px;
    height: 40px;
    border: 3px solid rgba(255, 255, 255, 0.1);
    border-radius: 50%;
    border-top-color: #4299e1;
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  @media (max-width: 768px) { 
    .dashboard {
      padding: 1rem;
    }
    
    .announcement-content {
      gap: 0.5rem;
    }
    
    .announcement-content p {
      font-size: 0.8rem;
    }

    .search-container {
      flex-direction: column;
    }

    .search-button {
      width: 100%;
      padding: 0.75rem;
    }

    .results-grid {
      grid-template-columns: 1fr;
    }

}
/* Mueve estos estilos FUERA del media query */
.error-bar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    background-color: rgba(220, 38, 38, 0.1); /* Rojo */
    border-left: 3px solid #dc2626; /* Rojo */
    border-radius: 0.375rem;
    padding: 0.75rem 1rem;
    margin-bottom: 1.5rem;
    transition: all 0.3s ease;
}

.error-content {
    display: flex;
    align-items: center;
    gap: 0.75rem;
}

.error-icon {
    color: #dc2626; /* Rojo */
    flex-shrink: 0;
}

.error-content p {
    margin: 0;
    font-size: 0.875rem;
    color: #fecaca; /* Rojo claro */
}

.error-close {
    background: none;
    border: none;
    color: #fca5a5; /* Rojo claro */
    cursor: pointer;
    padding: 0.25rem;
    border-radius: 0.25rem;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s ease;
}

.error-close:hover {
    color: #dc2626; /* Rojo intenso */
    background-color: rgba(220, 38, 38, 0.1);
}

@media (max-width: 768px) { 
    /* Mantén solo los ajustes responsive aquí */
    .error-bar {
        padding: 0.5rem 1rem;
    }
    
    .error-content p {
        font-size: 0.8rem;
    }
}
</style>