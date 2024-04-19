async function handleFormSubmit(event) {
    try {
        event.preventDefault()
        const response = await fetch('login/user', {
            method: 'POST',
            body: new FormData(loginForm),
        })

        const usernameInput = document.getElementById('username');
        const passwordInput = document.getElementById('password');
        const errorElement = document.getElementById('status');

        if (response.status === 200) {
            const response2 = await fetch('/account', {
                method: 'GET',
            })
            loginForm.reset();

            if (response2.status === 200) {
                window.location.href = '/account';
                return null
            }

            if (!response2.ok) {
                throw new Error(`Ошибка по адресу '/account', статус ошибки ${response2.status}`);
            }
            return null
        }

        if (response.status === 204) {
            errorElement.classList.add('text-danger');
            errorElement.textContent = 'Неверный логин или пароль';

            usernameInput.classList.add('is-invalid');
            passwordInput.classList.add('is-invalid');

            usernameInput.addEventListener('input', () => {
                errorElement.textContent = "";
                usernameInput.classList.remove('is-invalid');
                passwordInput.classList.remove('is-invalid');
            });

            passwordInput.addEventListener('input', () => {
                errorElement.textContent = "";
                usernameInput.classList.remove('is-invalid');
                passwordInput.classList.remove('is-invalid');
            });
            return null;
        }

        if (response.status === 401) {
            errorElement.classList.add('text-danger');
            errorElement.textContent = 'Неверный пароль';

            passwordInput.classList.add('is-invalid');

            passwordInput.addEventListener('input', () => {
                errorElement.textContent = "";
                usernameInput.classList.remove('is-invalid');
                passwordInput.classList.remove('is-invalid');
            });
            return null;
        }

        if (!response.ok) {
            throw new Error(`Ошибка по адресу 'login/user', статус ошибки ${response.status}`);
        }
    }
    catch (err) {
        console.log(err)
    }
}

function login() {
    const loginForm = document.getElementById('loginForm')
    loginForm.addEventListener('submit', handleFormSubmit)
}