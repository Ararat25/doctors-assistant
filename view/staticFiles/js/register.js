async function handleFormSubmit(event) {
    event.preventDefault()

    const usernameInput = document.getElementById('username');
    const errorElement = document.getElementById('status');
    const passwordInput = document.getElementById('password');
    const confirmPassword = document.getElementById('confirmPassword');

    const password = passwordInput.value;
    const passwordValid = password.length >= 8 &&
        /[A-Z]/.test(password) &&
        /[a-z]/.test(password) &&
        /[0-9]/.test(password) &&
        /[^A-Za-z0-9]/.test(password);

    if (!passwordValid) {
        passwordInput.classList.add('is-invalid');
        return;
    } else {
        passwordInput.classList.remove('is-invalid');
        passwordInput.classList.add('is-valid');
    }

    if (confirmPassword.value !== passwordInput.value) {
        confirmPassword.classList.add('is-invalid');
        return;
    } else {
        confirmPassword.classList.remove('is-invalid');
        confirmPassword.classList.add('is-valid');
    }

    const response = await fetch('register/user', {
        method: 'POST',
        body: new FormData(registerForm),
    })


    if (response.status === 409) {
        errorElement.classList.add('text-danger');
        errorElement.textContent = 'Аккаунт с такой почтой уже существует';

        usernameInput.classList.add('is-invalid');

        usernameInput.addEventListener('input', () => {
            errorElement.textContent = "";
            usernameInput.classList.remove('is-invalid');
        });
        return null;
    }

    if (!response.ok) {
        throw new Error(`Ошибка по адресу 'register/user', статус ошибки ${response.status}`);
    }

    if (response.status === 200) {
        registerForm.reset();
        window.location.href = '/login';
        return null
    }
}

function register() {
    const registerForm = document.getElementById('registerForm')
    registerForm.addEventListener('submit', handleFormSubmit)
}
