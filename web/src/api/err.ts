class ResponseError extends Error {
	public statusCode: number;
	public details: string;
	public name: string;

	constructor(message: string, statusCode: number, details: string, name: string) {
		super(message);
		this.statusCode = statusCode;
		this.details = details;
		this.name = name;
	}
}

export default ResponseError;