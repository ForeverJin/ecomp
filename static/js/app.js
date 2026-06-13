// Auto-dismiss flash alerts after 5 seconds
document.addEventListener('DOMContentLoaded', function() {
    document.querySelectorAll('.alert-dismissible').forEach(function(alert) {
        setTimeout(function() {
            if (alert.parentElement) {
                alert.style.transition = 'opacity 0.4s ease';
                alert.style.opacity = '0';
                setTimeout(function() { alert.remove(); }, 400);
            }
        }, 5000);
    });

    var searchInput = document.getElementById('searchInput');
    var tableBody = document.getElementById('compTableBody');
    if (searchInput && tableBody) {
        searchInput.addEventListener('input', function() {
            var keyword = searchInput.value.trim().toLowerCase();
            tableBody.querySelectorAll('tr').forEach(function(row) {
                var haystack = [
                    row.dataset.name,
                    row.dataset.model,
                    row.dataset.location,
                    row.dataset.package
                ].join(' ').toLowerCase();
                row.style.display = haystack.indexOf(keyword) === -1 ? 'none' : '';
            });
        });
    }
});

// Sidebar toggle (mobile)
function toggleSidebar() {
    var sidebar = document.getElementById('sidebar');
    var overlay = document.getElementById('sidebar-overlay');
    if (sidebar) {
        sidebar.classList.toggle('open');
    }
    if (overlay) {
        overlay.classList.toggle('active');
    }
}

function toggleAll(source) {
    document.querySelectorAll('.row-check').forEach(function(checkbox) {
        checkbox.checked = source.checked;
    });
}

// Close sidebar when clicking outside on mobile
document.addEventListener('click', function(e) {
    var sidebar = document.getElementById('sidebar');
    var overlay = document.getElementById('sidebar-overlay');
    if (sidebar && sidebar.classList.contains('open')) {
        if (!sidebar.contains(e.target) && !e.target.closest('.sidebar-toggle')) {
            sidebar.classList.remove('open');
            if (overlay) overlay.classList.remove('active');
        }
    }
});
