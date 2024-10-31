function toggleForms() {
    const registerForm = document.getElementById('registerForm');
    const loginForm = document.getElementById('loginForm');
    const toggleButton = document.getElementById('toggleButton');
    const labelText = document.getElementById('labelText');

    if (registerForm.style.display === "none") {
        registerForm.style.display = "block";
        loginForm.style.display = "none";
        toggleButton.innerText = "Switch to Login";
        labelText.innerText = "Wan't to create new account?";
    } else {
        registerForm.style.display = "none";
        loginForm.style.display = "block";
        toggleButton.innerText = "Switch to Register";
        labelText.innerText = "Don't have an account?";
    }
}
