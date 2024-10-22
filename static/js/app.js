document.getElementById('registerForm').addEventListener('submit', async (event) => {
    event.preventDefault();
    const username = document.getElementById('registerUsername').value;
    const password = document.getElementById('registerPassword').value;

    const response = await fetch('/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ username, password })
    });

    if (response.ok) {
        alert('Registration successful!');
        window.location.href = '/dashboard';
    } else {
        alert('Registration failed!');
    }
});

document.getElementById('loginForm').addEventListener('submit', async (event) => {
    event.preventDefault();
    const username = document.getElementById('loginUsername').value;
    const password = document.getElementById('loginPassword').value;

    const response = await fetch('/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ username, password })
    });

    if (response.ok) {
        const data = await response.json();
        const token = data.token;

        document.cookie = `jwtToken=${token}; path=/; Secure; SameSite=Strict`;

        alert('Login successful!');
        window.location.href = '/dashboard';
    } else {
        alert('Login failed!');
    }
});

async function accessDashboard() {
    const token = getCookie('jwtToken');
    
    const response = await fetch('/dashboard', {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });

    if (response.ok) {
        const dashboardData = await response.text();
        console.log(dashboardData);
    } else {
        alert('Failed to access dashboard!');
    }
}

function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
}

accessDashboard();