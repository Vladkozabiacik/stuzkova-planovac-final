// Register form handler
document.getElementById('registerForm').addEventListener('submit', async (event) => {
    event.preventDefault();
    const username = document.getElementById('registerUsername').value;
    const email = document.getElementById('registerEmail').value; // Extract email value
    const password = document.getElementById('registerPassword').value;

    try {
        console.log('Sending registration request for username:', username);
        
        const response = await fetch('http://127.0.0.1:8080/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            },
            credentials: 'include',
            body: JSON.stringify({ username, email, password }) // Include email in the request body
        });

        if (response.ok) {
            const data = await response.json();
            alert('Registration successful! Please login.');
            document.getElementById('registerForm').reset();
        } else {
            const errorData = await response.json().catch(() => response.text());
            const errorMessage = errorData.error || errorData;
            alert(`Registration failed: ${errorMessage}`);
        }
    } catch (error) {
        console.error('Registration error:', error);
        alert('Registration error: ' + error.message);
    }
});

// Login form handler
document.getElementById('loginForm').addEventListener('submit', async (event) => {
    event.preventDefault();
    const username = document.getElementById('loginUsername').value;
    const password = document.getElementById('loginPassword').value;

    try {
        console.log('Sending login request with username:', username);

        const response = await fetch('http://127.0.0.1:8080/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            },
            credentials: 'include',
            body: JSON.stringify({ username, password })
        });

        if (response.ok) {
            const data = await response.json();
            
            if (data.token) {
                // Store token in both cookie and localStorage for redundancy
                document.cookie = `jwtToken=${data.token}; path=/; SameSite=Strict; secure`;
                localStorage.setItem('jwtToken', data.token);
                
                window.location.href = '/dashboard';
            } else {
                throw new Error('No token received from server');
            }
        } else {
            const errorData = await response.json().catch(() => response.text());
            const errorMessage = errorData.error || errorData;
            console.error('Login failed:', errorMessage);
            alert(`Login failed: ${errorMessage}`);
        }
    } catch (error) {
        console.error('Login error:', error);
        alert('Login error: ' + error.message);
    }
});

// Dashboard access function
async function accessDashboard() {
    const token = getCookie('jwtToken') || localStorage.getItem('jwtToken');

    if (!token) {
        console.log('No token found, redirecting to login');
        window.location.href = '/';
        return;
    }

    try {
        const response = await fetch('http://127.0.0.1:8080/dashboard', {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Cache-Control': 'no-cache'
            },
            credentials: 'include'
        });

        if (response.ok) {
            const contentType = response.headers.get('content-type');
            if (contentType && contentType.includes('application/json')) {
                const data = await response.json();
                // Handle JSON response
                console.log('Dashboard data:', data);
            } else {
                // Handle HTML response
                const html = await response.text();
                document.documentElement.innerHTML = html;
            }
        } else if (response.status === 401) {
            console.log('Unauthorized, clearing tokens');
            // Clear both cookie and localStorage
            document.cookie = 'jwtToken=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
            localStorage.removeItem('jwtToken');
            window.location.href = '/';
        } else {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
    } catch (error) {
        console.error('Error accessing dashboard:', error);
        alert('An error occurred: ' + error.message);
    }
}

// Enhanced cookie getter
function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) {
        const cookieValue = parts.pop().split(';').shift();
        return cookieValue || null;
    }
    return null;
}

// Initialize dashboard access on page load
document.addEventListener('DOMContentLoaded', () => {
    const currentPath = window.location.pathname;
    if (currentPath === '/dashboard') {
        accessDashboard();
    }
});

// Add logout functionality
function logout() {
    document.cookie = 'jwtToken=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
    localStorage.removeItem('jwtToken');
    window.location.href = '/';
}
