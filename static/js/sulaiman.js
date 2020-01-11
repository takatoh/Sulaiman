$(() => {
  let app;
  let uploadDialog;
  let deleteDialog;
  let triggerFlag = false;

  const upload = (event) => {
    const ext = $("#upload-form input[name=file]").val().split(".").pop().toLowerCase();
    if ($.inArray(ext, ["jpg", "jpeg", "png", "gif"]) == -1) {
      alert("Unsupported format!\nSupport only jpg, png, gif.");
      clearUploadVals();
    } else {
      const fd = new FormData($("#upload-form").get(0));
      $.ajax({
        url: "/upload",
        type: "POST",
        data: fd,
        contentType: false,
        processData: false,
        dataType: "json"
      }).done((response) => {
        clearUploadVals();
        app.photoList.unshift(response.photo);
        if (response.deletePhotoId > 0) {
          app.photoList = app.photoList.filter((p) => {
            return p.id != response.deletePhotoId;
          });
        }
        uploadDialog.dialog("close");
      });
    }
  }

  const clearUploadVals = () => {
    $("#upload-form input[name=file]").val("");
    $("#upload-form input[name=key]").val("");
  }

  const deletePhoto = (event) => {
    const fd = new FormData($("#delete-form").get(0));
    $.ajax({
      url: "/delete",
      type: "DELETE",
      data: fd,
      contentType: false,
      processData: false,
      dataType: "json"
    }).done((response) => {
      if (response.status == "OK") {
        clearDeleteVals();
        app.photoList.some((v, i) => {
          if (v.id == response.photoId) {
            app.photoList.splice(i, 1);
          }
        });
        alert("Deleted: " + response.photoId);
      } else {
        clearDeleteVals();
        alert("Error! CAN'T delete: " + response.photoId);
      }
      deleteDialog.dialog("close");
    });
  }

  const clearDeleteVals = () => {
    $("#delete-form input[name=id]").val("");
    $("#delete-form input[name=key]").val("");
  }

  $.ajax({
    type: "GET",
    url: "/title",
    dataType: "text"
  }).done((response) => {
    $("head title").text(response)
    $("h1 a").text(response);
  });

  const nextUrl = $("#next-link").attr("href");
  $.ajax({
    type: "GET",
    url: nextUrl,
    dataType: "json"
  }).done((response) => {
    if (response.photos) {
      photoList = response.photos;
    } else {
      photoList = [];
    }
    app = new Vue({
      el: "#content",
      data: {
        photoList: photoList
      }
    });
    if (response.next) {
      $("#next-link").attr("href", response.next);
    } else {
      $("#next-link").remove();
    }
  });

  uploadDialog = $("#upload-dialog").dialog({
    autoOpen: false,
    modal: true,
    draggable: false,
    resizable: false,
    width: 350,
    buttons: {
      Upload: upload,
      Cancel: () => {
        clearUploadVals();
        uploadDialog.dialog("close");
      }
    }
  });

  $("#upload").button().on("click", () => {
    uploadDialog.dialog("open");
  });

  deleteDialog = $("#delete-dialog").dialog({
    autoOpen: false,
    modal: true,
    draggable: false,
    resizable: false,
    width: 300,
    buttons: {
      Delete: deletePhoto,
      Cancel: () => {
        clearDeleteVals();
        deleteDialog.dialog("close");
      }
    }
  });

  $("#delete").button().on("click", () => {
    deleteDialog.dialog("open");
  });

  $(window).on("load scroll", () => {
    const documentHeight = $(document).height();
    const scrollBottomPosition = $(window).height() + $(window).scrollTop();
    const triggerPoint = documentHeight - scrollBottomPosition;
    if (!triggerFlag && triggerPoint <= 50) {
      triggerFlag = true;
      const nextUrl = $("#next-link").attr("href");
      $.ajax({
        type: "GET",
        url: nextUrl,
        dataType: "json"
      }).done((response) => {
        if (response.photos) {
          response.photos.forEach((v) => { app.photoList.push(v); });
        }
        if (response.next) {
          $("#next-link").attr("href", response.next);
        } else {
          $("#next-link").remove();
        }
      });
    }
    if (triggerPoint > 50) {
      triggerFlag = false;
    }
  });
});
