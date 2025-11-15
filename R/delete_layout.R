#' Delete an existing project
#'
#' @param layout_name A string denoting the name of the pubSEM layout to
#'   delete.
#' @returns nothing
#' @export
delete_layout <- function(layout_name) {
    base_dir <- tools::R_user_dir("pubSEM", which = "data")
    layout_file <- file.path(base_dir, paste0(layout_name, ".json"))

    file.remove(layout_file)
    print(paste("Successfully deleted the following layout:", layout_name))
}
