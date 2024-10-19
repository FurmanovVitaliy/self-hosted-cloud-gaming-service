
export function validateEmail(email: string): boolean {
	const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
	return emailRegex.test(email);
}


export function validatePassword(password: string): boolean {
    // Регулярное выражение для пароля:
    // - Минимум 8 символов
    // - Хотя бы одна заглавная буква
    // - Хотя бы одна строчная буква
    // - Хотя бы одна цифра
    // - Хотя бы один специальный символ
    const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&\-#^])[A-Za-z\d@$!%*?&\-#^]{8,}$/;
    return passwordRegex.test(password);
  }

export function eventColorizeInput(e: Event) {
    let timeout: number | undefined;
    const target = e.target as HTMLInputElement;
    if (target && target.style) {
        target.style.color = "white";
    }
    if (timeout) {
        clearTimeout(timeout);
    }
    timeout = window.setTimeout(() => {
            target.style.color = target.validity.valid ? "white" : "rgba(237, 78, 78, 1)";
    },500);
}


export function hidePassword (e: Event){
    const target = e.target as HTMLButtonElement;
    const inputs = document.querySelectorAll('.input-box__input input[password]');
    if (inputs) {
        inputs.forEach((input) => {
            if (input.getAttribute('type') === 'password') {
                input.setAttribute('type', 'text');
                if (input.nextElementSibling) {
                  input.nextElementSibling.children[0].classList.remove('icon--eye-off');
                  input.nextElementSibling.children[0].classList.add('icon--eye-on');
                }
            } else {
                input.setAttribute('type', 'password');
                if (input.nextElementSibling) {
                    input.nextElementSibling.children[0].classList.remove('icon--eye-on');
                    input.nextElementSibling.children[0].classList.add('icon--eye-off');
                }
            }
        });
    }
   

  
  }

