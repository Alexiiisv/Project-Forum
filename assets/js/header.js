var btn = document.querySelector(".info");
var toshow = document.querySelector(".hide");
var test = true;

toshow.classList.add("off");

btn.addEventListener("mouseover", event => {
    test = false;
    toshow.classList.add("on");
    toshow.classList.remove("off");
});
btn.addEventListener("mouseout", event => {
    test = true;
    setTimeout(function() {
        if (test) {
            toshow.classList.add("off");
            toshow.classList.remove("on");
        }
    }, 2000);
});