// HTMX Debug Helper - Logs all HTMX events to console
(function() {
    'use strict';
    
    // Guard against double registration
    if (window.__htmxDebugLoaded) {
        return;
    }
    window.__htmxDebugLoaded = true;
    
    console.log('[HTMX Debug] Initializing debug helper...');
    
    // Log all HTMX events
    const events = [
        'htmx:beforeRequest',
        'htmx:afterRequest',
        'htmx:beforeSwap',
        'htmx:afterSwap',
        'htmx:afterSettle',
        'htmx:beforeProcessNode',
        'htmx:load',
        'htmx:responseError',
        'htmx:sendError',
        'htmx:targetError',
        'htmx:swapError'
    ];
    
    events.forEach(eventName => {
        document.addEventListener(eventName, function(evt) {
            console.log(`[HTMX Debug] ${eventName}:`, {
                target: evt.target?.id || evt.target?.tagName,
                detail: evt.detail,
                url: evt.detail?.xhr?.responseURL,
                status: evt.detail?.xhr?.status,
                response: evt.detail?.xhr?.response?.substring(0, 200)
            });
        });
    });
    
    // Log CSS load errors
    document.addEventListener('DOMContentLoaded', function() {
        const stylesheets = document.querySelectorAll('link[rel="stylesheet"]');
        console.log('[HTMX Debug] Loaded stylesheets:', Array.from(stylesheets).map(s => s.href));
        
        // Check for CSS load errors
        stylesheets.forEach(sheet => {
            sheet.addEventListener('error', function(e) {
                console.error('[HTMX Debug] Failed to load CSS:', sheet.href);
            });
        });
    });
    
    // Log Alpine.js initialization
    document.addEventListener('alpine:init', function() {
        console.log('[HTMX Debug] Alpine.js initialized');
    });
    
    console.log('[HTMX Debug] Debug helper active');
})();
