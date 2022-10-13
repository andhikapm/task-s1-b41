const realFileBtn = document.getElementById("input-blog-image");
const customBtn = document.getElementById("button-image");
const customTxt = document.getElementById("text-image");

customBtn.addEventListener("click", function() {
  realFileBtn.click();
});

realFileBtn.addEventListener("change", function() {
  if (realFileBtn.value) {
    customTxt.innerHTML = realFileBtn.value.match(
      /[\/\\]([\w\d\s\.\-\(\)]+)$/
    )[1];
  } else {
    customTxt.innerHTML = "";
  }
});


let dataBlog = []

function addBlog(event) {
    event.preventDefault()

    let title = document.getElementById("input-title").value
    let content = document.getElementById("input-content").value
    let image = document.getElementById("input-blog-image").files[0]

    image = URL.createObjectURL(image)
    console.log(image)

    let blog = {
        title,
        content,
        image,
        postAt: new Date(),
        author: "rangga alfa"
    }
    
    dataBlog.push(blog)
    console.log(dataBlog)

    localStorage.setItem("alpha" ,JSON.stringify(dataBlog))

    renderBlog()
}

function renderBlog() {
    document.getElementById("contents").innerHTML = ''

    for (let index = 0; index < dataBlog.length; index++) {
        console.log("test",dataBlog[index])

        document.getElementById("contents").innerHTML += `
        <div class="blog-list-item">
            <div class="blog-image">
                <img src="${dataBlog[index].image}">
            </div>
            <div class="blog-content">
                <div class="btn-group">
                    <button class="btn-edit">Edit Post</button>
                    <button class="btn-post">Post Blog</button>
                </div>
                <h1>
                    <a href="blog-detail.html" target="_blank">
                        ${dataBlog[index].title}
                    </a>
                </h1>
                <div class="detail-blog-content">
                    ${dataBlog[index].postAt} | ${dataBlog[index].author}
                </div>
                <p>
                    ${dataBlog[index].content}
                </p>
            </div>
        </div>
        `
    }
}