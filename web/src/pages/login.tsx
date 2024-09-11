import Page from '@/components/templates/page';


class LoginPage extends Page {
    private  html =(
        <h1>This is login page</h1>
    )
    public title: string= "Authentification";
    

    constructor() {
        super();
        this.insertTemplate(this.html);
    }
    
        async render() {
        return this.container;
    }
}
export default LoginPage;

