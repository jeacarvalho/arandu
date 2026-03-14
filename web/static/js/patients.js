// Patient management JavaScript
document.addEventListener('DOMContentLoaded', function() {
    // Form validation for new patient form
    const newPatientForm = document.getElementById('new-patient-form') || document.querySelector('form[action="/patients"]');
    
    if (newPatientForm) {
        newPatientForm.addEventListener('submit', function(e) {
            const nameInput = this.querySelector('input[name="name"]');
            const notesInput = this.querySelector('textarea[name="notes"]');
            let isValid = true;
            
            // Clear previous errors
            clearErrors(this);
            
            // Validate name
            if (!nameInput.value.trim()) {
                showError(nameInput, 'O nome do paciente é obrigatório.');
                isValid = false;
            } else if (nameInput.value.trim().length > 255) {
                showError(nameInput, 'O nome não pode exceder 255 caracteres.');
                isValid = false;
            }
            
            // Validate notes length
            if (notesInput && notesInput.value.length > 5000) {
                showError(notesInput, 'As observações não podem exceder 5000 caracteres.');
                isValid = false;
            }
            
            if (!isValid) {
                e.preventDefault();
                // Scroll to first error
                const firstError = this.querySelector('.form-input-error');
                if (firstError) {
                    firstError.scrollIntoView({ behavior: 'smooth', block: 'center' });
                    firstError.focus();
                }
            } else {
                // Show loading state
                const submitButton = this.querySelector('button[type="submit"]');
                if (submitButton) {
                    submitButton.disabled = true;
                    submitButton.innerHTML = '<span class="spinner"></span> Salvando...';
                }
            }
        });
        
        // Real-time validation
        const nameInput = newPatientForm.querySelector('input[name="name"]');
        if (nameInput) {
            nameInput.addEventListener('input', function() {
                if (this.value.trim().length > 255) {
                    showError(this, 'O nome não pode exceder 255 caracteres.');
                } else {
                    clearError(this);
                }
            });
        }
        
        const notesInput = newPatientForm.querySelector('textarea[name="notes"]');
        if (notesInput) {
            notesInput.addEventListener('input', function() {
                const charCount = this.value.length;
                const counter = document.getElementById('notes-counter') || createCounter(this);
                
                counter.textContent = `${charCount}/5000 caracteres`;
                
                if (charCount > 5000) {
                    counter.classList.add('text-red-600');
                    showError(this, 'Limite de 5000 caracteres excedido.');
                } else {
                    counter.classList.remove('text-red-600');
                    clearError(this);
                }
            });
        }
    }
    
    // Patient list interactions
    const patientCards = document.querySelectorAll('a[href^="/patient/"]');
    patientCards.forEach(card => {
        card.addEventListener('mouseenter', function() {
            this.style.transform = 'translateY(-2px)';
            this.style.transition = 'transform 0.2s ease';
        });
        
        card.addEventListener('mouseleave', function() {
            this.style.transform = 'translateY(0)';
        });
    });
    
    // Search functionality for patient list
    const searchInput = document.getElementById('patient-search');
    if (searchInput) {
        searchInput.addEventListener('input', function() {
            const searchTerm = this.value.toLowerCase();
            const patientItems = document.querySelectorAll('a[href^="/patient/"]');
            
            patientItems.forEach(item => {
                const patientName = item.querySelector('h3').textContent.toLowerCase();
                const patientNotes = item.querySelector('p.text-gray-600')?.textContent.toLowerCase() || '';
                
                if (patientName.includes(searchTerm) || patientNotes.includes(searchTerm)) {
                    item.style.display = 'block';
                } else {
                    item.style.display = 'none';
                }
            });
        });
    }
});

// Helper functions
function showError(input, message) {
    input.classList.add('form-input-error');
    
    let errorElement = input.nextElementSibling;
    if (!errorElement || !errorElement.classList.contains('form-error')) {
        errorElement = document.createElement('div');
        errorElement.className = 'form-error';
        input.parentNode.insertBefore(errorElement, input.nextSibling);
    }
    
    errorElement.textContent = message;
}

function clearError(input) {
    input.classList.remove('form-input-error');
    
    const errorElement = input.nextElementSibling;
    if (errorElement && errorElement.classList.contains('form-error')) {
        errorElement.remove();
    }
}

function clearErrors(form) {
    const errorInputs = form.querySelectorAll('.form-input-error');
    errorInputs.forEach(input => {
        input.classList.remove('form-input-error');
    });
    
    const errorMessages = form.querySelectorAll('.form-error');
    errorMessages.forEach(msg => msg.remove());
}

function createCounter(textarea) {
    const counter = document.createElement('div');
    counter.id = 'notes-counter';
    counter.className = 'text-sm text-gray-500 mt-1 text-right';
    counter.textContent = `0/5000 caracteres`;
    
    textarea.parentNode.appendChild(counter);
    return counter;
}

// Toast notification system
function showToast(message, type = 'success') {
    const toast = document.createElement('div');
    toast.className = `fixed top-4 right-4 z-50 px-4 py-3 rounded-lg shadow-lg transform transition-all duration-300 ${
        type === 'success' ? 'bg-green-100 text-green-800 border border-green-200' : 
        type === 'error' ? 'bg-red-100 text-red-800 border border-red-200' :
        'bg-blue-100 text-blue-800 border border-blue-200'
    }`;
    
    toast.innerHTML = `
        <div class="flex items-center">
            <span class="mr-2">${type === 'success' ? '✓' : type === 'error' ? '✗' : 'ℹ'}</span>
            <span>${message}</span>
        </div>
    `;
    
    document.body.appendChild(toast);
    
    // Auto-remove after 5 seconds
    setTimeout(() => {
        toast.style.opacity = '0';
        toast.style.transform = 'translateX(100%)';
        setTimeout(() => toast.remove(), 300);
    }, 5000);
}

// Export for use in other files
window.PatientUI = {
    showToast,
    showError,
    clearError
};