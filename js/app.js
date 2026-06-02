    const API_URL =
"http://localhost:6524";

const boton =
document.getElementById("btnLogin");
console.log(boton);

boton.addEventListener(
    "click",
function login(){
    alert("Boton funcionando");
});

async function login(){

    const usuario =
    document
    .getElementById("usuario")
    .value;

    const password =
    document
    .getElementById("password")
    .value;

    console.log(
        "Usuario:",
        usuario
    );

    console.log(
        "Password:",
        password
    );

    const auth =
    btoa(
        `${usuario}:${password}`
    );

    console.log(
        "Basic Auth:",
        auth
    );

    // backend futuro

    /*
    const respuesta =
    await fetch(
        `${API_URL}/login`,
        {

            method:"POST",

            headers:{

                Authorization:
                `Basic ${auth}`

            }

        }
    );
    */

}                                                                                                                                                              
 