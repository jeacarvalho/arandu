/**
 * HTMX Handlers - Global event handlers for HTMX
 * This file is loaded by both shell_layout.templ and layout.templ
 * 
 * Guard: window.__htmxHandlersLoaded prevents double registration
 */
(function() {
	'use strict';
	
	// Guard against double registration
	if (window.__htmxHandlersLoaded) {
		return;
	}
	window.__htmxHandlersLoaded = true;

	// Initialize Alpine stores for shared state
	document.addEventListener('alpine:init', () => {
		// Shell store - manages sidebar state for Shell layout
		Alpine.store('shell', {
			sidebarOpen: false,
			sidebarCollapsed: false,
			
			init() {
				// Restore sidebar collapsed state from localStorage
				const saved = localStorage.getItem('arandu-sidebar-collapsed');
				if (saved !== null) {
					this.sidebarCollapsed = saved === 'true';
				}
			},
			
			toggleSidebar() {
				this.sidebarOpen = !this.sidebarOpen;
			},
			
			toggleCollapse() {
				this.sidebarCollapsed = !this.sidebarCollapsed;
				localStorage.setItem('arandu-sidebar-collapsed', this.sidebarCollapsed);
			},
			
			closeSidebar() {
				this.sidebarOpen = false;
			}
		});
		
		// UI store - manages general UI state
		Alpine.store('ui', {
			sidebarOpen: false
		});
	});

	// HTMX Configuration
	htmx.config.defaultSwapDelay = 0;
	htmx.config.defaultSettleDelay = 0;
	htmx.config.scrollIntoViewOnBoost = false;

	// Global Error Handling with Toast
	function showErrorToast(message) {
		let toast = document.getElementById('error-toast');
		if (!toast) {
			toast = document.createElement('div');
			toast.id = 'error-toast';
			toast.className = 'fixed top-4 right-4 z-[100] max-w-sm bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg shadow-lg p-4 transform transition-all duration-300 translate-y-[-100%]';
			toast.setAttribute('role', 'alert');
			toast.setAttribute('aria-live', 'polite');
			document.body.appendChild(toast);
		}
		toast.innerHTML = `
			<div class="flex items-start gap-3">
				<i class="fas fa-exclamation-circle text-red-500 mt-0.5"></i>
				<div class="flex-1">
					<p class="text-sm text-red-800 dark:text-red-400 font-medium">Erro</p>
					<p class="text-sm text-red-600 dark:text-red-300 mt-1">${message}</p>
				</div>
				<button onclick="this.parentElement.parentElement.remove()" class="text-red-400 hover:text-red-600 dark:hover:text-red-200" aria-label="Fechar">
					<i class="fas fa-times"></i>
				</button>
			</div>
		`;
		toast.classList.remove('translate-y-[-100%]');
		toast.classList.add('translate-y-0');
		setTimeout(() => {
			toast.classList.add('translate-y-[-100%]');
			setTimeout(() => toast.remove(), 300);
		}, 5000);
	}

	// HTMX Error Events
	document.body.addEventListener('htmx:responseError', function(e) {
		const detail = e.detail;
		if (detail.xhr && detail.xhr.status >= 400) {
			let errorMessage = 'Erro ao processar requisição';
			try {
				const response = JSON.parse(detail.xhr.responseText);
				errorMessage = response.error || response.message || errorMessage;
			} catch {
				if (detail.xhr.responseText && detail.xhr.responseText.length < 200) {
					errorMessage = detail.xhr.responseText;
				}
			}
			showErrorToast(errorMessage);
		}
	});

	document.body.addEventListener('htmx:swapError', function(e) {
		showErrorToast('Erro ao atualizar conteúdo. Tente novamente.');
	});

// Focus management after HTMX swap
document.body.addEventListener('htmx:afterSwap', (e) => {
	// Close mobile sidebar on navigation using Alpine.store (official API)
	const shellStore = Alpine.store('shell');
	if (shellStore && typeof shellStore.closeSidebar === 'function' && window.innerWidth < 768) {
		shellStore.closeSidebar();
	}

	// Focus management
	const h1 = e.target.querySelector('h1, [data-autofocus]');
	if (h1) {
		h1.setAttribute('tabindex', '-1');
		h1.focus();
	}
});

// Handle head elements after HTMX navigation with hx-head="merge"
// This ensures CSS and other head elements are properly merged
document.body.addEventListener('htmx:beforeSwap', (e) => {
	// Store current scroll position before swap
	if (e.detail && e.detail.pathInfo) {
		sessionStorage.setItem('htmx-scroll-' + e.detail.pathInfo.requestPath, window.scrollY);
	}
});

// Restore scroll position after HTMX navigation
document.body.addEventListener('htmx:afterSettle', (e) => {
	// Restore scroll position if stored
	if (e.detail && e.detail.pathInfo) {
		const scrollY = sessionStorage.getItem('htmx-scroll-' + e.detail.pathInfo.requestPath);
		if (scrollY !== null) {
			window.scrollTo(0, parseInt(scrollY, 10));
			sessionStorage.removeItem('htmx-scroll-' + e.detail.pathInfo.requestPath);
		}
	}
});

// History restore handler - re-apply CSS classes after browser back/forward
// This fixes CSS breakage when navigating back to HTMX-saved snapshots
document.addEventListener('htmx:historyRestore', function(e) {
  console.log('[HTMX] History restore triggered');

  // Re-initialize Alpine.js on restored content
  if (typeof Alpine !== 'undefined') {
    document.querySelectorAll('[x-data]').forEach(el => {
      if (el._x_dataStack) {
        return;
      }
      Alpine.initTree(el);
    });
  }

  // Mark restoration complete
  document.body.classList.add('htmx-history-restored');
  setTimeout(() => {
    document.body.classList.remove('htmx-history-restored');
  }, 50);

  console.log('[HTMX] History restore completed');
});

  // Re-initialize Alpine.js on restored content
  if (typeof Alpine !== 'undefined') {
    document.querySelectorAll('[x-data]').forEach(el => {
      if (el._x_dataStack) {
        return;
      }
      Alpine.initTree(el);
    });
  }

  // Trigger reflow to ensure styles are recalculated
  document.body.style.display = 'none';
  document.body.offsetHeight; // Force reflow
  document.body.style.display = '';

  // Mark restoration complete
  document.body.classList.add('htmx-history-restored');
  setTimeout(() => {
    document.body.classList.remove('htmx-history-restored');
  }, 50);

  console.log('[HTMX] History restore completed, CSS refreshed');
});

	// Smooth page transitions
	document.addEventListener('DOMContentLoaded', function() {
		document.body.style.opacity = '0';
		requestAnimationFrame(() => {
			document.body.style.transition = 'opacity 0.3s ease';
			document.body.style.opacity = '1';
		});
	});

	// Expose showErrorToast globally for use in other scripts
	window.showErrorToast = showErrorToast;
})();