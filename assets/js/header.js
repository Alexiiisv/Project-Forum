var btn = document.querySelector(".info");
var toshow = document.querySelector(".hide");
btn.addEventListener("mouseover", event => {
    toshow.classList.add("on");
    toshow.classList.remove("off");
});
btn.addEventListener("mouseout", event => {
    setTimeout(function() {
        toshow.classList.add("off");
        toshow.classList.remove("on");
    }, 2000);
});