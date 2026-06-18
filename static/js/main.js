document.addEventListener('DOMContentLoaded', function () {
  var homeLink = document.querySelector('.nav-link[href="/"]');
  if (homeLink && window.location.pathname === '/') {
    homeLink.classList.add('active');
  }

  document.addEventListener('keydown', function (e) {
    if (e.key === 'Escape') {
      var offcanvas = document.querySelector('.offcanvas.show');
      if (offcanvas) {
        var instance = bootstrap.Offcanvas.getInstance(offcanvas);
        if (instance) instance.hide();
      }
    }

    if ((e.ctrlKey || e.metaKey) && e.key === 'f') {
      var searchInput = document.querySelector('input[name="q"]');
      if (searchInput && !e.target.closest('.offcanvas')) {
        e.preventDefault();
        searchInput.focus();
      }
    }

    if (e.key === 'ArrowLeft' || e.key === 'ArrowRight') {
      var links = document.querySelectorAll('.pagination .page-link');
      var prevLink, nextLink;
      links.forEach(function (link) {
        var text = link.textContent.trim().toLowerCase();
        if (text === 'previous' && !link.parentElement.classList.contains('disabled')) {
          prevLink = link;
        }
        if (text === 'next' && !link.parentElement.classList.contains('disabled')) {
          nextLink = link;
        }
      });

      if (e.key === 'ArrowLeft' && prevLink) {
        e.preventDefault();
        window.location.href = prevLink.getAttribute('href');
      }

      if (e.key === 'ArrowRight' && nextLink) {
        e.preventDefault();
        window.location.href = nextLink.getAttribute('href');
      }
    }
  });

  document.addEventListener('click', function (e) {
    var badge = e.target.closest('.filter-remove');
    if (!badge) return;
    e.preventDefault();

    var param = badge.getAttribute('data-filter-param');
    var value = badge.getAttribute('data-filter-value');
    var url = new URL(window.location.href);

    if (value) {
      var current = url.searchParams.getAll(param);
      url.searchParams.delete(param);
      current.forEach(function (v) {
        if (v !== value) {
          url.searchParams.append(param, v);
        }
      });
    } else {
      url.searchParams.delete(param);
    }

    window.location.href = url.toString();
  });
});
