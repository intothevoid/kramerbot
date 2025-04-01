document.addEventListener('DOMContentLoaded', () => {
    const tg = window.Telegram.WebApp;
    tg.ready(); // Inform Telegram the web app is ready

    // --- UI Elements ---
    const loadingDiv = document.getElementById('loading');
    const preferencesDiv = document.getElementById('preferences');
    const errorDiv = document.getElementById('error');

    const ozbGoodToggle = document.getElementById('ozbGood');
    const ozbSuperToggle = document.getElementById('ozbSuper');
    const amzDailyToggle = document.getElementById('amzDaily');
    const amzWeeklyToggle = document.getElementById('amzWeekly');

    const keywordsListUl = document.getElementById('keywordsList');
    const newKeywordInput = document.getElementById('newKeyword');
    const addKeywordBtn = document.getElementById('addKeywordBtn');

    const testBtn = document.getElementById('testBtn');

    // --- Functions ---
    function showError(message) {
        errorDiv.textContent = `Error: ${message}`; 
        errorDiv.style.display = 'block';
        loadingDiv.style.display = 'none';
        preferencesDiv.style.display = 'none';
    }

    function renderKeywords(keywords) {
        keywordsListUl.innerHTML = ''; // Clear existing list
        if (keywords && keywords.length > 0) {
            keywords.forEach(keyword => {
                const li = document.createElement('li');
                li.textContent = keyword;
                const removeBtn = document.createElement('button');
                removeBtn.textContent = 'Remove';
                removeBtn.onclick = () => handleRemoveKeyword(keyword);
                li.appendChild(removeBtn);
                keywordsListUl.appendChild(li);
            });
        } else {
            const li = document.createElement('li');
            li.textContent = 'No keywords added yet.';
            li.style.color = 'var(--tg-theme-hint-color)';
            keywordsListUl.appendChild(li);
        }
    }

    async function fetchPreferences() {
        loadingDiv.style.display = 'block';
        preferencesDiv.style.display = 'none';
        errorDiv.style.display = 'none';

        try {
            console.log("Fetching preferences from API...");
            // TODO: Replace with actual API endpoint
            const response = await fetch('/api/preferences', {
                 method: 'GET',
                 headers: {
                     // Pass user init data for authentication/identification on backend
                     'X-Telegram-Init-Data': tg.initData || 'dummy-init-data-for-testing'
                 }
             });

            console.log("Response status:", response.status);
            if (!response.ok) {
                const errorText = await response.text();
                console.error("API Error Response:", errorText);
                throw new Error(`HTTP error! status: ${response.status}, message: ${errorText}`);
            }
            
            const data = await response.json();
            console.log("Preferences data received:", data);

            // Populate UI
            ozbGoodToggle.checked = data.ozbGood;
            ozbSuperToggle.checked = data.ozbSuper;
            amzDailyToggle.checked = data.amzDaily;
            amzWeeklyToggle.checked = data.amzWeekly;
            renderKeywords(data.keywords);

            loadingDiv.style.display = 'none';
            preferencesDiv.style.display = 'block';
        } catch (err) {
            console.error("Failed to fetch preferences:", err);
            console.error("Error details:", err.message, err.stack);
            showError('Could not load your preferences. Please try again later.');
        }
    }

    async function updatePreference(key, value) {
        console.log(`Updating ${key} to ${value}`);
        try {
            // TODO: Replace with actual API endpoint
             const response = await fetch('/api/preferences', {
                 method: 'POST',
                 headers: {
                     'Content-Type': 'application/json',
                     'X-Telegram-Init-Data': tg.initData
                 },
                 body: JSON.stringify({ [key]: value })
             });

             if (!response.ok) {
                const errorData = await response.text(); 
                throw new Error(`Failed to update preference: ${errorData}`);
             }
             console.log(`Successfully updated ${key}`);
             // Optionally: Add visual feedback (e.g., a small checkmark)
        } catch (err) {
            console.error("Failed to update preference:", err);
            showError(`Failed to update ${key}. Please try reloading.`);
            // Revert UI change on failure? Or fetch prefs again?
            // For simplicity now, just show error.
        }
    }
    
    async function handleAddKeyword() {
        const keyword = newKeywordInput.value.trim().toLowerCase();
        if (!keyword) return;

        console.log(`Adding keyword: ${keyword}`);
        addKeywordBtn.disabled = true; // Prevent double clicks

        try {
            // TODO: Replace with actual API endpoint
            const response = await fetch('/api/keywords/add', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'X-Telegram-Init-Data': tg.initData
                },
                body: JSON.stringify({ keyword: keyword })
            });

            if (!response.ok) {
                const errorData = await response.text(); 
                throw new Error(`Failed to add keyword: ${errorData}`);
            }
            
            const updatedKeywords = await response.json(); // Expect backend to return the full list
            renderKeywords(updatedKeywords.keywords);
            newKeywordInput.value = ''; // Clear input
            console.log(`Successfully added keyword: ${keyword}`);
        } catch (err) {
            console.error("Failed to add keyword:", err);
            showError(err.message || 'Could not add keyword.');
        } finally {
            addKeywordBtn.disabled = false;
        }
    }

    async function handleRemoveKeyword(keyword) {
        console.log(`Removing keyword: ${keyword}`);
        // Optionally disable the specific button or show loading state

        try {
            // TODO: Replace with actual API endpoint
            const response = await fetch('/api/keywords/remove', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'X-Telegram-Init-Data': tg.initData
                },
                body: JSON.stringify({ keyword: keyword })
            });

             if (!response.ok) {
                const errorData = await response.text(); 
                throw new Error(`Failed to remove keyword: ${errorData}`);
             }

            const updatedKeywords = await response.json(); // Expect backend to return the full list
            renderKeywords(updatedKeywords.keywords);
            console.log(`Successfully removed keyword: ${keyword}`);
        } catch (err) {
            console.error("Failed to remove keyword:", err);
            showError(err.message || 'Could not remove keyword.');
        }
    }

    async function handleTestNotification() {
        console.log("Sending test notification...");
        testBtn.disabled = true;
        testBtn.textContent = 'Sending...';

        try {
            // TODO: Replace with actual API endpoint
            const response = await fetch('/api/test', {
                method: 'POST',
                 headers: {
                     'X-Telegram-Init-Data': tg.initData
                 }
            });
            
             if (!response.ok) {
                const errorData = await response.text(); 
                throw new Error(`Failed to send test notification: ${errorData}`);
             }

            console.log("Test notification request sent successfully.");
             // Maybe show a temporary success message instead of error
             errorDiv.textContent = 'Test notification sent!';
             errorDiv.style.color = 'green'; // Indicate success
             errorDiv.style.display = 'block';
             setTimeout(() => { errorDiv.style.display = 'none'; errorDiv.style.color = '#dc3545'; }, 3000);

        } catch (err) {
            console.error("Failed to send test notification:", err);
            showError(err.message || 'Could not send test notification.');
        } finally {
            testBtn.disabled = false;
            testBtn.textContent = 'Send Test Notification';
        }
    }

    // --- Event Listeners ---
    ozbGoodToggle.addEventListener('change', (e) => updatePreference('ozbGood', e.target.checked));
    ozbSuperToggle.addEventListener('change', (e) => updatePreference('ozbSuper', e.target.checked));
    amzDailyToggle.addEventListener('change', (e) => updatePreference('amzDaily', e.target.checked));
    amzWeeklyToggle.addEventListener('change', (e) => updatePreference('amzWeekly', e.target.checked));

    addKeywordBtn.addEventListener('click', handleAddKeyword);
    newKeywordInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            handleAddKeyword();
        }
    });

    testBtn.addEventListener('click', handleTestNotification);

    // --- Initialization ---
    fetchPreferences(); // Load initial data when the app loads

    // Apply Telegram theme parameters
    tg.expand(); // Expand the web app to full height
    // Optional: Set background color based on theme if needed explicitly
    // document.body.style.backgroundColor = tg.themeParams.bg_color || '#ffffff';
}); 