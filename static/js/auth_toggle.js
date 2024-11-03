function toggleForms() {
    const registerForm = document.getElementById('registerForm');
    const loginForm = document.getElementById('loginForm');
    
    // Toggle the visibility of the forms
    registerForm.classList.toggle('form-hidden');
    loginForm.classList.toggle('form-hidden');
}
