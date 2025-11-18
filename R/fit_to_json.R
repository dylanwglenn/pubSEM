# Writes a lavaan fit to temp.json
fit_to_json <- function(fit, standardized) {
    df_fit <- extract_lavaan_params(fit, standardized)

    base_dir <- tools::R_user_dir("pubSEM", which = "data")
    file_path <- file.path(base_dir, "temp.json")

    if (!dir.exists(file.path(base_dir))) {
        dir.create(file.path(base_dir), recursive = TRUE)
    }

    jsonlite::write_json(df_fit,
                         path = file_path,
                         pretty = TRUE
    )
}


# Get selected data from a lavaan fit
# returns a data frame
extract_lavaan_params <- function(fit, standardized) {
    df_paramTable <- lavaan::parametertable(fit)

    if (standardized) {
        df_estimates <- lavaan::standardizedsolution(fit)
        names(df_estimates)[names(df_estimates) == "est.std"] <- "est"
    } else {
        df_estimates <- lavaan::parameterestimates(fit)
    }
    #
    # select desired columns from fit data
    df_paramTable_filtered <- df_paramTable[, c("lhs",
                                                "op",
                                                "rhs",
                                                "user",
                                                "label",
                                                "group")]

    df_estimates_filtered <- df_estimates[, c("lhs",
                                              "op",
                                              "rhs",
                                              "est",
                                              "pvalue",
                                              "ci.lower",
                                              "ci.upper")]

    # merge the two tables into one
    df <- merge(df_paramTable_filtered, df_estimates_filtered)
    # change problematic names
    names(df)[names(df) == "ci.lower"] <- "ci_lower"
    names(df)[names(df) == "ci.upper"] <- "ci_upper"
    return(df)
}
