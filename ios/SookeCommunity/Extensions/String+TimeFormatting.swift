//
//  String+TimeFormatting.swift
//  SookeCommunity
//
//  Created by Kiefer Hay on 2026-03-29.
//

import Foundation

extension String {
    private static let timeInputFormatter: DateFormatter = {
        let f = DateFormatter()
        f.dateFormat = "HH:mm:ss"
        return f
    }()

    private static let timeOutputFormatter: DateFormatter = {
        let f = DateFormatter()
        f.dateFormat = "h:mm a"
        return f
    }()

    var formattedAsTime: String {
        guard let date = Self.timeInputFormatter.date(from: self) else {
            return self
        }
        return Self.timeOutputFormatter.string(from: date)
    }
}
