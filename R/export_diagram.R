
#' Export a pubSEM project to PDF
#'
#' @param project_name a string denoting the pubSEM project to export
#' @param filename a string specifying the name of the exported PDF
#' @param directory a string specifying the directory in which to save the
#'   exported PDF. Defaults to the current working directory.
#' @returns nothing
#' @export
export_diagram <- function(project_name, filename, directory = getwd()) {
    base_dir <- tools::R_user_dir("pubSEM", which = "data")
    file_path <- file.path(base_dir, paste(project_name, ".json"))

    if (string_end(filename, 4) != ".pdf") {
        filename <- paste0(filename,".pdf")
    }

    export_path <- file.path(directory, filename)

    if (Sys.info()['sysname'] == "Windows") {
        gui_exec_path <- system.file("bin", "sem_gui.exe", package = "pubSEM", mustWork = TRUE)
    } else {
        gui_exec_path <- system.file("bin", "sem_gui", package = "pubSEM", mustWork = TRUE)
    }

    # run the GUI executable
    system2(gui_exec_path,
            args = c(base_dir, project_name, "export", export_path)
    )
}

# get the last n characters from a string
string_end <- function(x, n) {
    substr(x, nchar(x)-n+1, nchar(x))
}
