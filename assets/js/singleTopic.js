var topic = document.querySelectorAll(".Topics");
var AccName = document.querySelectorAll(".AccName");
// console.log(AccName)
// topic.forEach(element => {
//     console.log("client " + element.clientHeight);

// });
function a() {
    console.log("c'est load")
        for (let i = 0; i < topic.length; i++) {
    AccName[i].style.heigth = "1000px";
    console.log(AccName[i].style.heigth)
}
}



window.onload = a;
console.log(AccName);