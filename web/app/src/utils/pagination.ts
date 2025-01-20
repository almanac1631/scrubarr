export interface Page {
    pageNumber: number | null,
    isSelected: boolean
}

export function getPageList(
    totalItems: number,
    pageSize: number,
    currentPage: number,
    maxPageButtons?: number  // Optional; if undefined, all pages are shown without collapsing
): Page[] {
    const totalPages = Math.ceil(totalItems / pageSize);
    const validatedCurrentPage = Math.min(Math.max(1, currentPage), totalPages);
    const pages: Page[] = [];

    // If maxPageButtons is undefined, display all pages without collapsing
    if (maxPageButtons === undefined) {
        for (let i = 1; i <= totalPages; i++) {
            pages.push({pageNumber: i, isSelected: i === validatedCurrentPage});
        }
        return pages;
    }

    // Helper function to add pages and placeholders
    const addPage = (pageNumber: number | null, isSelected: boolean = false) => {
        pages.push({pageNumber, isSelected});
    };

    // Always include the first page
    addPage(1, validatedCurrentPage === 1);

    // Add leading placeholder if needed and ensure pages near currentPage are always visible
    if (validatedCurrentPage > 3) {
        addPage(null); // Leading placeholder ("...")
    }

    // Main page range: includes `currentPage`, and one page before and after it
    const startPage = Math.max(2, validatedCurrentPage - 1); // Start from either page 2 or one before current
    const endPage = Math.min(totalPages - 1, validatedCurrentPage + 1); // End at one after current or second to last

    for (let i = startPage; i <= endPage; i++) {
        addPage(i, i === validatedCurrentPage);
    }

    // Add trailing placeholder if needed
    if (validatedCurrentPage < totalPages - 2) {
        addPage(null);  // Trailing placeholder ("...")
    }

    // Always include the last page
    if (totalPages > 1) {
        addPage(totalPages, validatedCurrentPage === totalPages);
    }

    return pages;
}
