// Registration Handler
document.getElementById('registerForm').addEventListener('submit', async (event) => {
    event.preventDefault();
    const username = document.getElementById('registerUsername').value;
    const password = document.getElementById('registerPassword').value;

    try {
        console.log('Sending registration request for username:', username);
        
        const response = await fetch('http://127.0.0.1:8080/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            },
            body: JSON.stringify({ username, password })
        });

        console.log('Registration response status:', response.status);

        if (response.ok) {
            const data = await response.json();
            console.log('Registration successful:', data);
            alert('Registration successful! Please login.');
            document.getElementById('registerForm').reset();
        } else {
            const errorText = await response.text();
            console.error('Registration failed:', errorText);
            alert(`Registration failed: ${errorText}`);
        }
    } catch (error) {
        console.error('Registration error:', error);
        alert('Registration error: ' + error.message);
    }
});

// Login Handler
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
            console.log('Login successful:', data);
            
            const token = data.token;
            document.cookie = `jwtToken=${token}; path=/; SameSite=Strict`;
            
            alert('Login successful!');
            window.location.href = '/dashboard';
        } else {
            const errorText = await response.text();
            console.error('Login failed:', errorText);
            alert(`Login failed: ${errorText}`);
        }
    } catch (error) {
        console.error('Login error:', error);
        alert('Login error: ' + error.message);
    }
});

// Access Dashboard Function
async function accessDashboard() {
    const token = getCookie('jwtToken');

    try {
        const response = await fetch('/dashboard', {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Cache-Control': 'no-cache'
            }
        });

        if (response.ok) {
            window.location.href = '/dashboard';
        } else if (response.status === 401) {
            document.cookie = 'jwtToken=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
        } else {
            alert('Failed to access dashboard!');
        }
    } catch (error) {
        console.error('Error accessing dashboard:', error);
        alert('An error occurred: ' + error.message);
    }
}


function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
}

document.addEventListener('DOMContentLoaded', () => {
    accessDashboard();
});
