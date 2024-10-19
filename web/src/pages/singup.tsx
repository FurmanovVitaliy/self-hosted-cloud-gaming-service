import Page from "@/components/templates/page";
import authApi from "@/api/auth";
import ResponseError from "@/api/err";
import { eventColorizeInput, validateEmail, validatePassword, hidePassword } from "@/common/input-validations";
import eventBus from "@/common/event-bus";
import constants from "@/common/constants";

type singupResponseData = { id: string; email: string; username: string };

class SingUpPage extends Page {
	private html = (
		<null>
			<section class='auth'>
				<a class='link-button auth__close-button' amepad-focus-element href={constants.routes.home} inner-link animated-upd>
					<span>
						<i class='icon icon--go-close'></i>
					</span>
				</a>
				<div class='auth-container'>
					<form action='' animated-upd>
						<h2>SingUp</h2>
						<div class='auth-container__input input-box'>
							<label>Email</label>
							<input email type='email' required placeholder='Email'></input>
						</div>
						<div class='auth-container__input input-box'>
							<label>Username</label>
							<input username type='text' required placeholder='Username'></input>
						</div>
						<div class='auth-container__input input-box'>
							<label>Password</label>
							<div class='input-box__input'>
								<input password type='password' required autocomplete placeholder='Passworld'></input>
								<button type='button'>
									<i class='icon icon--eye-off'></i>
								</button>
							</div>
							<span></span>
						</div>
						<div class='auth-container__input input-box'>
							<label>Confirm Password</label>
							<div class='input-box__input'>
								<input password type='password' required autocomplete placeholder='Passworld'></input>
								<button type='button'>
									<i class='icon icon--eye-off'></i>
								</button>
							</div>
						</div>
						<button class='auth-container__button button button--main' type='submit'>
							SingUp
						</button>
						<p class='auth-form-change'>
							Already have an account?{" "}
							<a href={constants.routes.login} inner-link>
								Login
							</a>
						</p>
					</form>
				</div>
			</section>
		</null>
	);
	public title: string = "Authentification | Pixel Cloud";
	//declare main fields
	constructor() {
		super();
		//add content to page
		this.loadTemplate(this.html);
		//local events
		this.addEventListener("input[email]", "input", eventColorizeInput);
		this.addEventListener('button[type="submit"]', "click", this.handleSubmit);
		this.addEventListener(".input-box__input button", "click", hidePassword, true);
		//local functions
		//connectedCallback
	}
	//! all appeal to html documents as this.container
	private async handleSubmit(e: Event) {
		const emailInput = document.querySelector("input[email]") as HTMLInputElement;
		const username = document.querySelector("input[username]") as HTMLInputElement;
		const passwordInput = document.querySelectorAll("input[password]") as NodeListOf<HTMLInputElement>;
		e.preventDefault();

		if (!validateEmail(emailInput.value)) {
			alert("Please provide a valid email.");
			return;
		}
		if (!username.value) {
			alert("Please provide a username.");
			return;
		}

		if (!validatePassword(passwordInput[0].value)) {
			alert("Please provide a valid password.");
			return;
		}
		if (passwordInput[0].value !== passwordInput[1].value) {
			alert("Passwords do not match.");
			return;
		}
		const submitButton = document.querySelector('button[type="submit"]') as HTMLButtonElement;
		if (submitButton) {
			submitButton.disabled = true;
			submitButton.textContent = "Logging in...";
		}
		try {
			const responseData = (await authApi.signup(
				emailInput.value,
				username.value,
				passwordInput[0].value,
			)) as singupResponseData;
			localStorage.setItem("userId", responseData.id);
			localStorage.setItem("username", responseData.username);
			eventBus.notify("nav", constants.routes.login);
		} catch (error) {
			if (error instanceof ResponseError) {
				alert(error.details);
			} else {
				alert("An error occurred. Please try again later.");
			}
		} finally {
			if (submitButton) {
				submitButton.disabled = false;
				submitButton.textContent = "Log in";
			}
		}
	}

	async render() {
		return this.container;
	}
}
export default SingUpPage;
