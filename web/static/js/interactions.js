// Arandu Interactions - Micro-interactions and UI enhancements

document.addEventListener('DOMContentLoaded', function() {
    // Initialize all interactions
    initFormInteractions();
    initCardInteractions();
    initButtonInteractions();
    initTooltips();
    initLoadingStates();
    initToastNotifications();
    initHTMXFocusManagement();
});

// Form interactions
function initFormInteractions() {
    const formControls = document.querySelectorAll('.form-control');
    
    formControls.forEach(control => {
        // Add focus/blur effects
        control.addEventListener('focus', function() {
            this.parentElement.classList.add('focused');
        });
        
        control.addEventListener('blur', function() {
            this.parentElement.classList.remove('focused');
            validateField(this);
        });
        
        // Real-time validation for text inputs
        if (control.type === 'text' || control.type === 'textarea') {
            control.addEventListener('input', function() {
                validateField(this);
                updateCharacterCounter(this);
            });
        }
    });
    
    // Form submission animations
    const forms = document.querySelectorAll('form');
    forms.forEach(form => {
        form.addEventListener('submit', function(e) {
            const submitBtn = this.querySelector('button[type="submit"]');
            if (submitBtn) {
                submitBtn.classList.add('loading');
                submitBtn.innerHTML = '<i class="fas fa-spinner fa-spin btn-icon"></i> Processando...';
                submitBtn.disabled = true;
            }
        });
    });
}

// Field validation
function validateField(field) {
    const wrapper = field.parentElement;
    const errorElement = wrapper.nextElementSibling;
    
    // Clear previous states
    wrapper.classList.remove('valid', 'invalid');
    
    if (field.validity.valid) {
        if (field.value.trim() !== '') {
            wrapper.classList.add('valid');
        }
    } else {
        wrapper.classList.add('invalid');
        if (errorElement && errorElement.classList.contains('form-error')) {
            errorElement.style.display = 'flex';
        }
    }
}

// Character counter
function updateCharacterCounter(field) {
    const counter = field.parentElement.querySelector('.form-counter');
    if (!counter) return;
    
    const maxLength = field.maxLength || 5000;
    const currentLength = field.value.length;
    const remaining = maxLength - currentLength;
    
    counter.textContent = `${currentLength}/${maxLength}`;
    
    // Update styling based on remaining characters
    counter.classList.remove('warning', 'error');
    
    if (remaining < 100 && remaining >= 0) {
        counter.classList.add('warning');
    } else if (remaining < 0) {
        counter.classList.add('error');
    }
}

// Card interactions
function initCardInteractions() {
    const cards = document.querySelectorAll('.card-hover');
    
    cards.forEach(card => {
        // Hover effect
        card.addEventListener('mouseenter', function() {
            this.style.transform = 'translateY(-8px)';
            this.style.boxShadow = 'var(--shadow-xl)';
        });
        
        card.addEventListener('mouseleave', function() {
            this.style.transform = 'translateY(0)';
            this.style.boxShadow = 'var(--shadow-md)';
        });
        
        // Click ripple effect
        card.addEventListener('click', function(e) {
            if (this.tagName === 'A' || this.tagName === 'BUTTON') return;
            
            const ripple = document.createElement('span');
            const rect = this.getBoundingClientRect();
            const size = Math.max(rect.width, rect.height);
            const x = e.clientX - rect.left - size / 2;
            const y = e.clientY - rect.top - size / 2;
            
            ripple.style.cssText = `
                position: absolute;
                border-radius: 50%;
                background: rgba(45, 125, 229, 0.1);
                transform: scale(0);
                animation: ripple-animation 0.6s linear;
                width: ${size}px;
                height: ${size}px;
                top: ${y}px;
                left: ${x}px;
                pointer-events: none;
            `;
            
            this.style.position = 'relative';
            this.style.overflow = 'hidden';
            this.appendChild(ripple);
            
            setTimeout(() => ripple.remove(), 600);
        });
    });
}

// Button interactions
function initButtonInteractions() {
    const buttons = document.querySelectorAll('.btn');
    
    buttons.forEach(btn => {
        // Click animation
        btn.addEventListener('click', function(e) {
            if (this.classList.contains('loading')) return;
            
            // Ripple effect
            const ripple = document.createElement('span');
            const rect = this.getBoundingClientRect();
            const size = Math.max(rect.width, rect.height);
            const x = e.clientX - rect.left - size / 2;
            const y = e.clientY - rect.top - size / 2;
            
            ripple.style.cssText = `
                position: absolute;
                border-radius: 50%;
                background: rgba(255, 255, 255, 0.3);
                transform: scale(0);
                animation: ripple-animation 0.6s linear;
                width: ${size}px;
                height: ${size}px;
                top: ${y}px;
                left: ${x}px;
                pointer-events: none;
            `;
            
            this.appendChild(ripple);
            setTimeout(() => ripple.remove(), 600);
            
            // Icon animation
            const icon = this.querySelector('.btn-icon');
            if (icon) {
                icon.style.transform = 'scale(1.2) rotate(5deg)';
                setTimeout(() => {
                    icon.style.transform = 'scale(1) rotate(0)';
                }, 200);
            }
        });
        
        // Hover effect
        btn.addEventListener('mouseenter', function() {
            this.style.transform = 'translateY(-2px)';
        });
        
        btn.addEventListener('mouseleave', function() {
            this.style.transform = 'translateY(0)';
        });
    });
}

// Tooltips
function initTooltips() {
    const tooltipElements = document.querySelectorAll('[data-tooltip]');
    
    tooltipElements.forEach(element => {
        element.addEventListener('mouseenter', function() {
            const tooltip = document.createElement('div');
            tooltip.className = 'tooltip-popup';
            tooltip.textContent = this.dataset.tooltip;
            
            const rect = this.getBoundingClientRect();
            tooltip.style.cssText = `
                position: fixed;
                background: var(--neutral-900);
                color: white;
                padding: var(--space-xs) var(--space-sm);
                border-radius: var(--radius-sm);
                font-size: 0.75rem;
                z-index: 10000;
                white-space: nowrap;
                pointer-events: none;
                transform: translateX(-50%);
                top: ${rect.top - 40}px;
                left: ${rect.left + rect.width / 2}px;
            `;
            
            document.body.appendChild(tooltip);
            this._tooltip = tooltip;
        });
        
        element.addEventListener('mouseleave', function() {
            if (this._tooltip) {
                this._tooltip.remove();
                this._tooltip = null;
            }
        });
    });
}

// Loading states
function initLoadingStates() {
    // Add CSS for ripple animation
    const style = document.createElement('style');
    style.textContent = `
        @keyframes ripple-animation {
            to {
                transform: scale(4);
                opacity: 0;
            }
        }
        
        .loading-spinner {
            display: inline-block;
            width: 1.25rem;
            height: 1.25rem;
            border: 2px solid var(--neutral-200);
            border-top-color: var(--primary-500);
            border-radius: 50%;
            animation: spin 1s linear infinite;
            margin-right: var(--space-sm);
        }
        
        @keyframes spin {
            to { transform: rotate(360deg); }
        }
    `;
    document.head.appendChild(style);
}

// Toast notifications
function initToastNotifications() {
    window.showToast = function(message, type = 'info') {
        const toast = document.createElement('div');
        toast.className = `toast toast-${type}`;
        toast.innerHTML = `
            <div class="toast-content">
                <i class="fas fa-${getToastIcon(type)}"></i>
                <span>${message}</span>
            </div>
            <button class="toast-close">
                <i class="fas fa-times"></i>
            </button>
        `;
        
        toast.style.cssText = `
            position: fixed;
            top: 20px;
            right: 20px;
            background: ${getToastColor(type)};
            color: white;
            padding: var(--space-md) var(--space-lg);
            border-radius: var(--radius-md);
            box-shadow: var(--shadow-lg);
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: var(--space-md);
            z-index: 10000;
            animation: slideIn 0.3s ease-out;
            max-width: 400px;
        `;
        
        document.body.appendChild(toast);
        
        // Close button
        const closeBtn = toast.querySelector('.toast-close');
        closeBtn.addEventListener('click', () => toast.remove());
        
        // Auto-remove after 5 seconds
        setTimeout(() => {
            if (toast.parentNode) {
                toast.style.animation = 'fadeOut 0.3s ease-out';
                setTimeout(() => toast.remove(), 300);
            }
        }, 5000);
    };
    
    function getToastIcon(type) {
        const icons = {
            success: 'check-circle',
            error: 'exclamation-circle',
            warning: 'exclamation-triangle',
            info: 'info-circle'
        };
        return icons[type] || 'info-circle';
    }
    
    function getToastColor(type) {
        const colors = {
            success: 'var(--success)',
            error: 'var(--error)',
            warning: 'var(--warning)',
            info: 'var(--info)'
        };
        return colors[type] || 'var(--info)';
    }
    
    // Add fadeOut animation
    const fadeOutStyle = document.createElement('style');
    fadeOutStyle.textContent = `
        @keyframes fadeOut {
            from { opacity: 1; transform: translateX(0); }
            to { opacity: 0; transform: translateX(100%); }
        }
    `;
    document.head.appendChild(fadeOutStyle);
}

// Form success/error handling
window.handleFormSuccess = function(form, message = 'Operação realizada com sucesso!') {
    const submitBtn = form.querySelector('button[type="submit"]');
    if (submitBtn) {
        submitBtn.classList.remove('loading');
        submitBtn.disabled = false;
        submitBtn.innerHTML = '<i class="fas fa-check btn-icon"></i> Concluído';
        submitBtn.classList.add('btn-success');
        
        setTimeout(() => {
            submitBtn.innerHTML = submitBtn.dataset.originalText || 'Salvar';
            submitBtn.classList.remove('btn-success');
        }, 2000);
    }
    
    showToast(message, 'success');
};

window.handleFormError = function(form, message = 'Ocorreu um erro. Tente novamente.') {
    const submitBtn = form.querySelector('button[type="submit"]');
    if (submitBtn) {
        submitBtn.classList.remove('loading');
        submitBtn.disabled = false;
        submitBtn.innerHTML = submitBtn.dataset.originalText || 'Tentar Novamente';
    }
    
    showToast(message, 'error');
    
    // Shake form animation
    form.style.animation = 'shake 0.5s ease-in-out';
    setTimeout(() => {
        form.style.animation = '';
    }, 500);
};

// Add shake animation
const shakeStyle = document.createElement('style');
shakeStyle.textContent = `
    @keyframes shake {
        0%, 100% { transform: translateX(0); }
        10%, 30%, 50%, 70%, 90% { transform: translateX(-5px); }
        20%, 40%, 60%, 80% { transform: translateX(5px); }
    }
`;
document.head.appendChild(shakeStyle);

// HTMX Focus Management - Alta 3
function initHTMXFocusManagement() {
    const focusStyle = document.createElement('style');
    focusStyle.textContent = `
        @keyframes focusErrorPulse {
            0%, 100% { box-shadow: 0 0 0 0 transparent; }
            50% { box-shadow: 0 0 0 3px rgba(220, 38, 38, 0.3); }
        }
        .focus-error {
            animation: focusErrorPulse 0.5s ease-in-out 2;
        }
    `;
    document.head.appendChild(focusStyle);
    
    document.addEventListener('htmx:afterSwap', function(event) {
        handleFocusPostSwap(event.detail);
    });
    
    document.addEventListener('htmx:responseError', function(event) {
        handleErrorFocus(event.detail);
    });
}

function handleFocusPostSwap(detail) {
    const target = detail.target;
    if (!target) return;
    
    setTimeout(() => {
        const errorInput = target.querySelector('[aria-invalid="true"], .invalid input, .invalid textarea, input[aria-describedby*="error"], textarea[aria-describedby*="error"]');
        
        if (errorInput) {
            errorInput.focus();
            errorInput.scrollIntoView({ behavior: 'smooth', block: 'center' });
            return;
        }
        
        const errorAlert = target.querySelector('[role="alert"], .form-error, .error-message');
        if (errorAlert) {
            const form = target.closest('form');
            if (form) {
                const firstInput = form.querySelector('input:not([type="hidden"]), textarea, select');
                if (firstInput) {
                    firstInput.focus();
                    firstInput.scrollIntoView({ behavior: 'smooth', block: 'center' });
                }
            }
            return;
        }
        
        const newItem = target.querySelector('input, textarea, select');
        if (newItem && !newItem.matches('[readonly], [disabled]')) {
            newItem.focus();
        }
    }, 50);
}

function handleErrorFocus(detail) {
    const target = detail.target;
    if (!target) return;
    
    setTimeout(() => {
        const firstInput = target.querySelector('input:not([type="hidden"]), textarea, select');
        if (firstInput) {
            firstInput.focus();
            firstInput.classList.add('focus-error');
            setTimeout(() => firstInput.classList.remove('focus-error'), 2000);
        }
    }, 50);
}